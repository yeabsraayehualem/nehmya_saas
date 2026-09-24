package tenants

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrInvalidSubscriptionDates = errors.New("subscription end date must be after its start date")

type TenantService struct {
	repository              Repository
	provisioner             DatabaseProvisioner
	defaultSubscriptionDays int
}

func NewService(repository Repository, provisioner DatabaseProvisioner, defaultSubscriptionDays int) *TenantService {
	if defaultSubscriptionDays <= 0 {
		defaultSubscriptionDays = 30
	}
	return &TenantService{repository: repository, provisioner: provisioner, defaultSubscriptionDays: defaultSubscriptionDays}
}

func (s *TenantService) Register(name string, ownerID uint) (*Tenant, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("tenant name is required")
	}
	start := time.Now().UTC()
	end := start.AddDate(0, 0, s.defaultSubscriptionDays)

	randomID := make([]byte, 10)
	if _, err := rand.Read(randomID); err != nil {
		return nil, fmt.Errorf("generate tenant database name: %w", err)
	}
	databaseName := "tenant_" + hex.EncodeToString(randomID)
	if err := s.provisioner.CreateDatabase(databaseName); err != nil {
		return nil, fmt.Errorf("create tenant database: %w", err)
	}

	tenant := &Tenant{
		Name: name, OwnerUserID: ownerID, DatabaseName: databaseName, Status: StatusActive,
		Subscription: Subscription{StartsAt: start, EndsAt: end},
	}
	if err := s.repository.Create(tenant); err != nil {
		if cleanupErr := s.provisioner.DropDatabase(databaseName); cleanupErr != nil {
			return nil, fmt.Errorf("save tenant: %v (database cleanup also failed: %w)", err, cleanupErr)
		}
		return nil, fmt.Errorf("save tenant: %w", err)
	}
	return tenant, nil
}

func (s *TenantService) ListByOwner(ownerID uint) ([]Tenant, error) {
	return s.repository.ListByOwner(ownerID)
}

func (s *TenantService) SetSubscriptionDates(tenantID uint, startsAt, endsAt time.Time) error {
	startsAt, endsAt = startsAt.UTC(), endsAt.UTC()
	if !endsAt.After(startsAt) {
		return ErrInvalidSubscriptionDates
	}
	return s.repository.UpdateSubscriptionDates(tenantID, startsAt, endsAt)
}
