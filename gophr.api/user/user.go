package user

import (
	"github.com/go-playground/validator/v10"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

type User struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty" validate:"required"`
	Email    string `json:"email,omitempty" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,gte=8,lte=130"`
	Note string
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
	user.ID = GenerateID()
	return user, nil
}
