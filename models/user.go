package models

type UserCreation struct {
	Password    string `validate:"required"`
	Email       string `validate:"required,email"`
	DisplayName string `validate:"required"`
	Bio         string
}

type Login struct {
	Password string `validate:"required"`
	Email    string `validate:"required,email"`
}
