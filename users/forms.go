package users

import ("errors"
	"fmt"
)


type CreateUserDTO struct {
	Name string
	Email string
	Phone string
	Password string
	ConfirmPassword string
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