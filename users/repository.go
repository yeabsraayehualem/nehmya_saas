package users

import (
	"errors"

	"gorm.io/gorm"

	"github.com/yeabsraayehualem/nehmya_saas/utils"
)

// ErrInvalidCredentials is returned when a login attempt fails.
var ErrInvalidCredentials = errors.New("invalid email/phone or password")

// Repository defines the persistence operations for users.
type Repository interface {
	Create(user *User) error
	GetByID(id uint) (*User, error)
	GetByPhone(phone string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	// Authenticate verifies the login (email or phone) against the stored
	// bcrypt hash and returns the user on success.
	Authenticate(login, password string) (*User, error)
}

// UserRepository is the gorm-backed Repository implementation.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository returns a gorm-backed user repository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id uint) (*User, error) {
	var user User

	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByPhone(phone string) (*User, error) {
	var user User

	if err := r.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	var user User

	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(user *User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Authenticate(login, password string) (*User, error) {
	var user User

	if err := r.db.Where("email = ? OR phone = ?", login, login).First(&user).Error; err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := utils.CompareHash(user.Password, password); err != nil {
		return nil, ErrInvalidCredentials
	}

	return &user, nil
}
