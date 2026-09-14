package users

import (
	"nehmya/cmd/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type UserHandler struct {
	service UserService
	store *sessions.CookieStore
}

func NewUserHandler(service UserService, store *sessions.CookieStore) *UserHandler{
	return &UserHandler{
		service: service,
		store: store,
	}
}



func (h UserHandler) CreateUser(c *gin.Context) {
	var form CreateUserDTO

	if err := c.ShouldBindJSON(&form); err != nil{
		utils.Error(c,http.StatusBadRequest,err.Error())
		return
	}

	if err := form.Valid(); err != nil{
		utils.Error(c,http.StatusBadRequest,err.Error())
		return
	}

	user := User{
		Name: form.Name,
		Email: form.Email,
		Password: form.Password,
		Phone: form.Phone,
	}

	h.service.CreateUser(user)
}

