package users

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"

	"github.com/yeabsraayehualem/nehmya_saas/utils"
)

type UserHandler struct {
	service *UserService
	store   *sessions.CookieStore
}

// NewUserHandler wires a user handler to its service and session store.
func NewUserHandler(service *UserService, store *sessions.CookieStore) *UserHandler {
	return &UserHandler{
		service: service,
		store:   store,
	}
}

// Authenticate logs a user in and stores the user ID in the session cookie.
func (h UserHandler) Authenticate(c *gin.Context) {
	var form LoginRequest

	if err := c.ShouldBindJSON(&form); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := form.Validate(); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	login := form.Phone
	if form.Email != "" {
		login = form.Email
	}

	user, err := h.service.Authenticate(login, form.Password)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	session, err := h.store.Get(c.Request, "session")
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to read session")
		return
	}
	session.Values["user_id"] = user.ID
	if err := session.Save(c.Request, c.Writer); err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to save session")
		return
	}

	utils.Success(c, http.StatusOK, "Authentication Successful!", user)
}

// CreateUser registers a new user.
func (h UserHandler) CreateUser(c *gin.Context) {
	var form CreateUserDTO

	if err := c.ShouldBindJSON(&form); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := form.Valid(); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	password, err := utils.HashText(form.Password)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	user := User{
		Name:     form.Name,
		Email:    form.Email,
		Password: password,
		Phone:    form.Phone,
	}

	if err := h.service.CreateUser(&user); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, "User created successfully!", user)
}

// GetMe returns the currently authenticated user from the session.
func (h UserHandler) GetMe(c *gin.Context) {
	session, err := h.store.Get(c.Request, "session")
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to read session")
		return
	}

	raw, ok := session.Values["user_id"]
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "not authenticated")
		return
	}

	id, ok := raw.(uint)
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "not authenticated")
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusUnauthorized, "not authenticated")
			return
		}
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Current user", user)
}

// RequireAuth is middleware that rejects requests without a valid session.
func (h UserHandler) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := h.store.Get(c.Request, "session")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to read session"})
			return
		}

		raw, ok := session.Values["user_id"]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "not authenticated"})
			return
		}

		id, ok := raw.(uint)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "not authenticated"})
			return
		}

		c.Set("user_id", id)
		c.Next()
	}
}
