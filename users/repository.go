package users

type Repository interface {
	Create(user *User) error
	GetByID(id uint) (*User, error)
	GetByPhone(phone string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
}