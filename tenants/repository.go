package tenants

import (
	"time"

	"gorm.io/gorm"
)

// Repository defines persistence operations for tenant registration.
type Repository interface {
	Create(tenant *Tenant) error
	ListByOwner(ownerID uint) ([]Tenant, error)
	UpdateSubscriptionDates(tenantID uint, startsAt, endsAt time.Time) error
}

type TenantRepository struct{ db *gorm.DB }

func NewTenantRepository(db *gorm.DB) *TenantRepository { return &TenantRepository{db: db} }

func (r *TenantRepository) Create(tenant *Tenant) error {
	return r.db.Transaction(func(tx *gorm.DB) error { return tx.Create(tenant).Error })
}

func (r *TenantRepository) ListByOwner(ownerID uint) ([]Tenant, error) {
	var tenants []Tenant
	err := r.db.Preload("Subscription").Where("owner_user_id = ?", ownerID).Order("created_at DESC").Find(&tenants).Error
	return tenants, err
}

func (r *TenantRepository) UpdateSubscriptionDates(tenantID uint, startsAt, endsAt time.Time) error {
	result := r.db.Model(&Subscription{}).Where("tenant_id = ?", tenantID).
		Updates(map[string]any{"starts_at": startsAt, "ends_at": endsAt})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
