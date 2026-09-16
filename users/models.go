package users

import "fmt"

const (
	RoleSysAdmin = "sys_admin"
	RoleManager  = "manager"
)

// User is the built-in account model. It is auto-migrated by the nehmya
// server on startup.
type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Phone    string `gorm:"uniqueIndex;not null" json:"phone"`
	Password string `gorm:"not null" json:"-"` // never serialized
	Active   bool   `gorm:"not null;default:true" json:"active"`
	Role     string `gorm:"not null;default:'manager'" json:"role"`
}

func (u *User) ValidateRole() error {
	if u.Role != RoleSysAdmin && u.Role != RoleManager {
		return fmt.Errorf("invalid role: %s", u.Role)
	}

	return nil
}
