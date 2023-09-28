package user

import (
	"ardaa/domain"
	"ardaa/internal/token"
	"ardaa/internal/validator"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type iUserService struct {
	repo *iUserRepo
}

// Login implements domain.UserService.
func (service *iUserService) Login(user *domain.LoginUser) (*map[string]interface{}, error) {
	// validate the user
	err := validator.Validate(user)
	if err != nil {
		return nil, err
	}

	ok := service.repo.EmailExists(user.Email)
	if !ok {
		return nil, errors.New("invalid credentials")
	}

	db_user, _ := service.repo.GetByEmail(user.Email)

	password_err := bcrypt.CompareHashAndPassword([]byte(db_user.Password), []byte(user.Password))
	if password_err != nil {
		return nil, errors.New("invalid credentials")
	}

	token_str, error := token.NewToken(db_user.UUID, time.Now().Add(24*time.Hour), user.Ip)
	if error != nil {
		return nil, error
	}

	res := map[string]interface{}{
		"token": token_str,
		"user":  db_user.AuthUser(),
	}

	return &res, nil
}

// Logout implements domain.UserService.
func (*iUserService) Logout(user_id string) error {
	return token.DeleteToken(user_id)
}

// Me implements domain.UserService.
func (service *iUserService) Me(id string) (*domain.AuthUser, error) {
	user, err := service.repo.Me(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Register implements domain.UserService.
func (service *iUserService) Register(user domain.RegisterUser) (*map[string]interface{}, error) {
	// validate the user
	validationErr := validator.Validate(user)
	if validationErr != nil {
		return nil, validationErr
	}

	// hash the Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	validUser := domain.User{
		Name:     user.Name,
		Email:    user.Email,
		Gender:   user.Gender,
		Password: string(hashedPassword),
		Username: user.Username,
		UUID:     uuid.New().String(),
	}

	validUser.Role = "regular"

	// check if the user already exists
	emailExists := service.repo.EmailExists(validUser.Email)
	if emailExists {
		return nil, errors.New("email already exists")
	}

	// check if the username already exists
	usernameExists := service.repo.UsernameExists(validUser.Username)
	if usernameExists {
		return nil, errors.New("username already exists")
	}

	// create the user
	err = service.repo.Create(&validUser)
	if err != nil {
		panic(err)
	}

	// returning the user with no password
	auth_token, err := token.NewToken(validUser.UUID, time.Now().Add(24*time.Hour), user.Ip)
	if err != nil {
		return nil, err
	}

	authUser := validUser.AuthUser()

	res := map[string]interface{}{
		"token": auth_token,
		"user":  authUser,
	}

	return &res, nil
}

// is user already logged in
func (*iUserService) IsLoggedIn(ip string) bool {
	record := token.GetToken(ip)

	return record != nil
}

// UpdateProfile implements domain.UserService.
func (*iUserService) UpdateProfile() {
	panic("unimplemented")
}

// NewUserService returns a new instance of domain.UserService.
func NewUserService(repo *iUserRepo) *iUserService {
	return &iUserService{
		repo,
	}
}

var _ domain.UserService = &iUserService{}
