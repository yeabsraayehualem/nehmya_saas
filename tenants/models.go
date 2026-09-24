package tenants

import "time"

const StatusActive = "active"

// Tenant is a customer organization registered by a signed-in account.
type Tenant struct {
	ID           uint         `gorm:"primaryKey" json:"id"`
	Name         string       `gorm:"not null" json:"name"`
	OwnerUserID  uint         `gorm:"not null;index" json:"owner_user_id"`
	DatabaseName string       `gorm:"uniqueIndex;not null" json:"database_name"`
	Status       string       `gorm:"not null;default:active" json:"status"`
	Subscription Subscription `gorm:"foreignKey:TenantID" json:"subscription"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// Subscription stores the service period attached to a tenant.
type Subscription struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TenantID  uint      `gorm:"uniqueIndex;not null" json:"tenant_id"`
	StartsAt  time.Time `gorm:"not null" json:"starts_at"`
	EndsAt    time.Time `gorm:"not null;index" json:"ends_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
