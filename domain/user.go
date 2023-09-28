package domain

import (
	"time"

	"gorm.io/gorm"
)

// User Waa user domain-ka matalaya system user
type DomainUser struct {
	UpdatedAt   *time.Time
	CreatedAt   *time.Time
	UUID        string
	Email       string
	Password    string
	Name        string
	Role        string
	Bio         string
	AvatarURL   string
	Username    string
	Gender      string
	IsVerified  bool
	IsSuspended bool
	IsFirstTime bool
}

// DB model of the user
type User struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Password    string         `gorm:"not null" json:"password"`
	Gender      string         `gorm:"type:enum('male', 'female');default:male" json:"gender"`
	Bio         string         `gorm:"size:250" json:"bio"`
	Username    string         `gorm:"unique;not null" json:"username"`
	Email       string         `gorm:"unique;not null" json:"email"`
	UUID        string         `gorm:"primaryKey;autoIncrement:false" json:"uuid"`
	Name        string         `gorm:"not null" json:"name"`
	Role        string         `gorm:"type:enum('admin', 'moderator', 'regular');default:regular" json:"role"`
	AvatarURL   string         `json:"avatar_url"`
	ID          uint           `gorm:"primaryKey"`
	IsVerified  bool           `gorm:"default:false" json:"is_verified"`
	IsSuspended bool           `gorm:"default:false" json:"is_suspended"`
	IsFirstTime bool           `gorm:"default:true" json:"is_first_time"`
}

// AuthUser specific for returning to the user like PublicUser
type AuthUser struct {
	UpdatedAt   *time.Time
	CreatedAt   *time.Time
	UUID        string
	Email       string
	Name        string
	Role        string
	Bio         string
	AvatarURL   string
	Username    string
	Gender      string
	IsVerified  bool
	IsSuspended bool
	IsFirstTime bool
}

// RegisterUser gaar u ah marka la register gareynayo user-ka
type RegisterUser struct {
	Name      string `validate:"required,validateNotBlank,min=5"`
	Username  string `validate:"required,validateNotBlank,validateAlphaWithSpecialChars"`
	Email     string `validate:"required,email"`
	AvatarURL string `validate:"omitempty,uri"`
	Gender    string `validate:"required,oneof=male Male Female female"` // because oneof is case-sensitive
	Password  string `validate:"required,validateNotBlank,min=8,validateAlphaWithSpecialChars"`
	Ip        []byte
}

// LoginUser gaar u ah marka la login gareynayo user-ka
type LoginUser struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,validateNotBlank,min=8,validateAlphaWithSpecialChars"`
	Ip       []byte
}

// UserRepository Database interface-ka that maps the use of the user struct
type UserRepository interface {
	Create(*User) error
	Get(string) (User, error)
	Update(string, *User) error
	Delete(string) error
	GetBy()    // waxan u samaynaynaa filter structs uu la xidhiidhi karo si uusan u batch search gareynin database-ka dhan
	GetAllBy() // kanna sida GetBy oo lakin helaya all matching users
	EmailExists(string) bool
	UsernameExists(string) bool
	GetByEmail(string) (*User, error)
	Me(string) (*AuthUser, error)
}

// UserService -ka la xidhiidhaya user Repository-ga
type UserService interface {
	Register(RegisterUser) (*map[string]interface{}, error)
	Login(*LoginUser) (*map[string]interface{}, error)
	Logout(string) error
	UpdateProfile()
	Me(string) (*AuthUser, error)
	IsLoggedIn(string) bool
}

func (u *User) AuthUser() *AuthUser {
	return &AuthUser{
		UpdatedAt:   &u.UpdatedAt,
		CreatedAt:   &u.CreatedAt,
		UUID:        u.UUID,
		Email:       u.Email,
		Name:        u.Name,
		Role:        u.Role,
		Bio:         u.Bio,
		AvatarURL:   u.AvatarURL,
		Username:    u.Username,
		Gender:      u.Gender,
		IsVerified:  u.IsVerified,
		IsSuspended: u.IsSuspended,
		IsFirstTime: u.IsFirstTime,
	}
}
