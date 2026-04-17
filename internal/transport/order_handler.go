package transport

import (
	"net/http"
	"pharmacy/internal/models"
	"pharmacy/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler(service *services.OrderService) OrderHandler {
	return OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) RegisterRoutes(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.POST("/:id/orders", h.CreateOrder)
	}

	orders := r.Group("/orders")
	{
		orders.GET("/:id", h.GetOrder)
		orders.PATCH("/:id/status", h.UpdateStatus)
	}
}

func (h *OrderHandler) CreateOrder(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var req models.OrderCreateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	order, err := h.service.CreateOrder(uint(userID), req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetOrder(ctx *gin.Context) {
	orderID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	order, err := h.service.GetOrder(uint(orderID))

	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "заказ не найден" {
			status = http.StatusNotFound
		}
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id is not a number"})
		return
	}

	var req models.OrderUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateStatus(uint(orderID), req); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "заказ не найден" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	order, err := h.service.GetOrder(uint(orderID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "order status updated"})
		return
	}

	c.JSON(http.StatusOK, order)
}
