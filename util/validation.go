package util

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type Validation struct {
	validate *validator.Validate
}

func (v Validation) ValidateStruct(input any) error {
	// returns nil or ValidationErrors ( []FieldError )
	err := v.validate.Struct(input)
	var msgs string
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			msgs += fmt.Sprintf("Field %s needs to be %s  ", err.Field(), err.Tag())
		}
		// from here you can create your own error messages in whatever language you wish
		return errors.New(msgs)
	}
	return nil
}

func NewValidation() *Validation {
	validate := validator.New()
	return &Validation{validate: validate}
}

func HashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Print(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}
