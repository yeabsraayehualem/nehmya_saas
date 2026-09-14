package users

import "fmt"

const (
	RoleSysAdmin = "sys_admin"
	RoleManager  = "manager"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Phone    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Active   bool   `gorm:"not null;default:true"`
	Role     string `gorm:"not null;default:'manager'"`
}

func (u *User) ValidateRole() error {
	if u.Role != RoleSysAdmin && u.Role != RoleManager {
		return fmt.Errorf("invalid role: %s", u.Role)
	}

	return nil
}