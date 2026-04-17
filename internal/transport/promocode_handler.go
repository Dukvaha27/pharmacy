package transport

import (
	"net/http"
	"pharmacy/internal/models"
	"pharmacy/internal/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PromocodeHandler struct {
	service *services.PromocodeService
}

func NewPromocodeHandler(service *services.PromocodeService) PromocodeHandler {
	return PromocodeHandler{
		service: service,
	}
}

func (h *PromocodeHandler) RegisterRoutes(r *gin.Engine) {
	promocodes := r.Group("/promocodes")
	{
		promocodes.GET("", h.GetAll)
		promocodes.POST("", h.Create)
		promocodes.POST("/validate", h.Validate)
		promocodes.PATCH("/:id", h.Update)
		promocodes.DELETE("/:id", h.Delete)
	}
}

func (h *PromocodeHandler) GetAll(c *gin.Context) {
	promocodes, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, promocodes)
}

func (h *PromocodeHandler) Create(c *gin.Context) {
	var req models.PromocodeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	promocode, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, promocode)
}

func (h *PromocodeHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "promocode id is not a number"})
		return
	}

	var req models.PromocodeUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	promocode, err := h.service.Update(uint(id), &req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "промокод не найден" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, promocode)
}

func (h *PromocodeHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "promocode id is not a number"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "промокод не найден" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "promocode deleted"})
}

func (h *PromocodeHandler) Validate(c *gin.Context) {
	var req models.PromocodeCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.OrderAmount < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_amount must be >= 0"})
		return
	}

	promocode, err := h.service.GetByCode(req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !promocode.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "промокод неактивен"})
		return
	}

	now := time.Now()
	if !promocode.ValidFrom.IsZero() && now.Before(promocode.ValidFrom) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "промокод ещё не действует"})
		return
	}
	if !promocode.ValidTo.IsZero() && now.After(promocode.ValidTo) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "срок действия промокода истёк"})
		return
	}
	if promocode.MaxUses > 0 && promocode.UsedCount >= promocode.MaxUses {
		c.JSON(http.StatusBadRequest, gin.H{"error": "достигнут лимит использований промокода"})
		return
	}

	orderAmount := uint64(req.OrderAmount)
	var discount uint64

	switch promocode.DiscountType {
	case "percent":
		discount = orderAmount * promocode.DiscountValue / 100
	case "fixed":
		discount = promocode.DiscountValue
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "неизвестный тип скидки"})
		return
	}

	if discount > orderAmount {
		discount = orderAmount
	}

	c.JSON(http.StatusOK, gin.H{
		"code":           promocode.Code,
		"is_valid":       true,
		"discount_total": discount,
		"final_price":    orderAmount - discount,
	})
}
