package users

import (
	"errors"
	"fmt"
)

type CreateUserDTO struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (c CreateUserDTO) Valid() error {
	fields := map[string]string{
		"Name":     c.Name,
		"Email":    c.Email,
		"Phone":    c.Phone,
		"Password": c.Password,
	}

	for field, value := range fields {
		if value == "" {
			return fmt.Errorf("%s field is required", field)
		}
	}

	if c.Password != c.ConfirmPassword {
		return errors.New("passwords do not match")
	}

	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (l LoginRequest) Validate() error {
	if l.Email == "" && l.Phone == "" {
		return errors.New("Phone or Email is required!")
	}

	if l.Password == "" {
		return errors.New("Password field is required!")
	}

	return nil
}
