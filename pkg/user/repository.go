package user

import (
	"ardaa/domain"

	"gorm.io/gorm"
)

type iUserRepo struct {
	db *gorm.DB
}

// Create implements domain.UserRepository.
func (repo *iUserRepo) Create(user *domain.User) error {
	return repo.db.Create(user).Error
}

// Delete implements domain.UserRepository.
func (*iUserRepo) Delete(string) error {
	panic("unimplemented")
}

// Get implements domain.UserRepository.
func (*iUserRepo) Get(string) (domain.User, error) {
	panic("unimplemented")
}

// GetAllBy implements domain.UserRepository.
func (*iUserRepo) GetAllBy() {
	panic("unimplemented")
}

// GetBy implements domain.UserRepository.
func (*iUserRepo) GetBy() {
	panic("unimplemented")
}

// Update implements domain.UserRepository.
func (*iUserRepo) Update(string, *domain.User) error {
	panic("unimplemented")
}

// EmailExists implements domain.UserRepository.
func (repo *iUserRepo) EmailExists(email string) bool {
	return repo.db.Where("email = ?", email).First(&domain.User{}).RowsAffected > 0
}

// UsernameExists implements domain.UserRepository.
func (repo *iUserRepo) UsernameExists(username string) bool {
	return repo.db.Where("username = ?", username).First(&domain.User{}).RowsAffected > 0
}

// GetByEmail implements domain.UserRepository.
func (repo *iUserRepo) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	repo.db.Where("email = ?", email).First(&domain.User{}).Scan(&user)
	return &user, nil
}

// Me implements domain.UserRepository. returns the user with the given id
func (repo *iUserRepo) Me(id string) (*domain.AuthUser, error) {
	var user domain.User
	repo.db.Where("uuid = ?", id).First(&domain.User{}).Scan(&user)

	return user.AuthUser(), nil
}

func NewUserRepo(db *gorm.DB) *iUserRepo {
	return &iUserRepo{
		db: db,
	}
}

var _ domain.UserRepository = &iUserRepo{}
