package user

import (
	"github.com/go-playground/validator/v10"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var validate = validator.New()

type User struct {
	UserID   string `json:"userId,omitempty" gorm:"user_id"`
	Username string `json:"username,omitempty" validate:"required" gorm:"username"`
	Email    string `json:"email,omitempty" validate:"required,email" gorm:"email"`
	Password string `json:"password,omitempty" validate:"required,gte=8,lte=130" gorm:"password"`

	// Base
	ID        uint `json:"id,omitempty"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time `sql:"index"`
}

func GenerateID() string {
	guid := xid.New()
	return guid.String()
}

func NewUser(username, email, password string) (*User, error) {
	user := &User{
		Username: username,
		Email:    email,
		Password: password,
	}
	err := validate.Struct(user)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)
	user.UserID = GenerateID()
	return user, nil
}
