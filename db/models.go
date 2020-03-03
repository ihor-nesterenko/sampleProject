package db

import validation "github.com/go-ozzo/ozzo-validation"

type User struct {
	Login    string
	Password string
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Login, validation.Required),
		validation.Field(&u.Password, validation.Required),
	)
}
