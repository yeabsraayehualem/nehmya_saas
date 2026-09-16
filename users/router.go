package users

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

// Routes mounts the built-in user endpoints on the given group:
//
//	POST /user/create  – register
//	POST /user/login   – authenticate
//	GET  /user/me      – current session user (protected)
func Routes(r gin.IRouter, db *gorm.DB, store *sessions.CookieStore) {
	repository := NewUserRepository(db)
	service := NewService(repository)
	handler := NewUserHandler(service, store)

	users := r.Group("/user")

	users.POST("/create", handler.CreateUser)
	users.POST("/login", handler.Authenticate)
	users.GET("/me", handler.RequireAuth(), handler.GetMe)
}
