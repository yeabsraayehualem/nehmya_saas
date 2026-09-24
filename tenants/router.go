package tenants

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/yeabsraayehualem/nehmya_saas/users"
	"gorm.io/gorm"
)

// Routes mounts authenticated tenant registration and listing endpoints.
func Routes(r gin.IRouter, db *gorm.DB, store *sessions.CookieStore, provisioner DatabaseProvisioner, defaultSubscriptionDays int) {
	handler := NewHandler(NewService(NewTenantRepository(db), provisioner, defaultSubscriptionDays))
	tenants := r.Group("/tenant", users.SessionAuth(store))
	tenants.POST("/register", handler.Register)
	tenants.GET("/mine", handler.ListMine)
	tenants.PUT("/:id/subscription", requireSysAdmin(db), handler.SetSubscriptionDates)
}

func requireSysAdmin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := c.Get("user_id")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "not authenticated"})
			return
		}
		var user users.User
		if err := db.Select("id", "role", "active").First(&user, id).Error; err != nil || !user.Active {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "not authenticated"})
			return
		}
		if user.Role != users.RoleSysAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "system administrator role required"})
			return
		}
		c.Next()
	}
}
