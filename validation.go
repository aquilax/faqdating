package main

import (
	"net/http"
	"net/mail"
)

type validationErrors []string

type Validation struct {
	ln *Language
}

func NewValidation(ln *Language) *Validation {
	return &Validation{ln}
}

func validEmail(email string) (string, error) {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", err
	}
	return addr.Address, nil
}

func (v *Validation) validateUserRegister(r *http.Request) (string, string, validationErrors) {
	password := r.FormValue("password")
	errors := validationErrors{}
	email, err := validEmail(r.FormValue("email"))
	if err != nil {
		errors = append(errors, v.ln.Lang("Please enter valid email address"))
	}
	if password != r.FormValue("password_confirm") {
		errors = append(errors, v.ln.Lang("Passwords don't match"))
	}
	if len(password) < 3 {
		errors = append(errors, v.ln.Lang("Passwords is too short"))
	}
	return email, password, errors
}
