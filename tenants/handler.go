package tenants

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yeabsraayehualem/nehmya_saas/utils"
)

type TenantHandler struct{ service *TenantService }

func NewHandler(service *TenantService) *TenantHandler { return &TenantHandler{service: service} }

type registerTenantRequest struct {
	Name string `json:"name" binding:"required"`
}

type setSubscriptionDatesRequest struct {
	StartsAt time.Time `json:"starts_at" binding:"required"`
	EndsAt   time.Time `json:"ends_at" binding:"required"`
}

func (h *TenantHandler) Register(c *gin.Context) {
	var request registerTenantRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	ownerID, ok := c.Get("user_id")
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "not authenticated")
		return
	}
	id, ok := ownerID.(uint)
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "not authenticated")
		return
	}
	tenant, err := h.service.Register(request.Name, id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusCreated, "Tenant registered and database provisioned", tenant)
}

func (h *TenantHandler) ListMine(c *gin.Context) {
	ownerID, ok := c.Get("user_id")
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "not authenticated")
		return
	}
	id, ok := ownerID.(uint)
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "not authenticated")
		return
	}
	tenants, err := h.service.ListByOwner(id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, "Your tenants", tenants)
}

func (h *TenantHandler) SetSubscriptionDates(c *gin.Context) {
	var request setSubscriptionDatesRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		utils.Error(c, http.StatusBadRequest, "invalid tenant id")
		return
	}
	if err := h.service.SetSubscriptionDates(uint(id), request.StartsAt, request.EndsAt); err != nil {
		if err == ErrInvalidSubscriptionDates {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.Error(c, http.StatusNotFound, "tenant subscription not found")
		return
	}
	utils.Success(c, http.StatusOK, "Subscription dates updated", nil)
}
