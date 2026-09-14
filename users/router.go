package users

import (

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)


func Routes(r *gin.Engine, db *gorm.DB,store *sessions.CookieStore){
	repository := NewUserRepository(db)
	service := NewService(*repository)
	handler := NewUserHandler(*service,store)



	users := r.Group("/user")

	users.POST("/create", handler.CreateUser)

	
	
}