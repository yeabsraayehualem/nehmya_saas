package tenants

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DatabaseProvisioner creates a separate PostgreSQL database for each tenant.
type DatabaseProvisioner interface {
	CreateDatabase(name string) error
	DropDatabase(name string) error
}

type PostgresProvisioner struct {
	host, port, user, password, adminDatabase string
}

func NewPostgresProvisioner(host, port, user, password, adminDatabase string) *PostgresProvisioner {
	if adminDatabase == "" {
		adminDatabase = "postgres"
	}
	return &PostgresProvisioner{host: host, port: port, user: user, password: password, adminDatabase: adminDatabase}
}

var databaseNamePattern = regexp.MustCompile(`^tenant_[a-f0-9]{20}$`)

func (p *PostgresProvisioner) openAdmin() (*gorm.DB, error) {
	if p.host == "" || p.port == "" || p.user == "" || p.adminDatabase == "" {
		return nil, fmt.Errorf("database provisioning requires DB_HOST, DB_PORT, DB_USER, and DB_ADMIN_DATABASE")
	}
	dsn := url.URL{Scheme: "postgres", User: url.UserPassword(p.user, p.password), Host: net.JoinHostPort(p.host, p.port), Path: "/" + p.adminDatabase}
	query := url.Values{}
	query.Set("sslmode", "disable")
	dsn.RawQuery = query.Encode()
	return gorm.Open(postgres.Open(dsn.String()), &gorm.Config{})
}

func (p *PostgresProvisioner) CreateDatabase(name string) error {
	if !databaseNamePattern.MatchString(name) {
		return fmt.Errorf("invalid tenant database name")
	}
	db, err := p.openAdmin()
	if err != nil {
		return fmt.Errorf("connect to PostgreSQL admin database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	return db.Exec(`CREATE DATABASE "` + strings.ReplaceAll(name, `"`, `""`) + `"`).Error
}

func (p *PostgresProvisioner) DropDatabase(name string) error {
	if !databaseNamePattern.MatchString(name) {
		return fmt.Errorf("invalid tenant database name")
	}
	db, err := p.openAdmin()
	if err != nil {
		return fmt.Errorf("connect to PostgreSQL admin database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()
	return db.Exec(`DROP DATABASE IF EXISTS "` + strings.ReplaceAll(name, `"`, `""`) + `"`).Error
}
