package server

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/yeabsraayehualem/nehmya_saas/cmd/utils"
	"github.com/yeabsraayehualem/nehmya_saas/tenants"
	"github.com/yeabsraayehualem/nehmya_saas/users"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var migrationModels = []interface{}{
	&users.User{},
	&tenants.Tenant{},
	&tenants.Subscription{},
}

type Server struct {
	PORT string
	DB   *gorm.DB
}

func NewServer() *Server {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on process environment")
	}

	constring := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(constring), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	autoMigrate(db)
	utils.InitSessions()

	port := ":" + os.Getenv("PORT")

	return &Server{
		PORT: port,
		DB:   db,
	}
}

func autoMigrate(db *gorm.DB) {
	for _, model := range migrationModels {
		if err := db.AutoMigrate(model); err != nil {
			log.Fatalf("failed to migrate: %v", err)
		}
	}
}

func (s *Server) Run() {
	r := gin.Default()

	users.Routes(r, s.DB, utils.CookieStore)
	provisioner := tenants.NewPostgresProvisioner(
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_ADMIN_DATABASE"),
	)
	defaultSubscriptionDays, err := strconv.Atoi(os.Getenv("DEFAULT_SUBSCRIPTION_DAYS"))
	if err != nil || defaultSubscriptionDays <= 0 {
		defaultSubscriptionDays = 30
	}
	tenants.Routes(r, s.DB, utils.CookieStore, provisioner, defaultSubscriptionDays)

	log.Println("Server is running on port " + s.PORT)

	if err := r.Run(s.PORT); err != nil {
		log.Fatal(err)
	}
}
