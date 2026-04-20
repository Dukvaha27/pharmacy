package transport

import (
	"net/http"
	"pharmacy/internal/models"
	"pharmacy/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service      *services.PaymentService
	orderService *services.OrderService
}

func NewPaymentHandler(
	service *services.PaymentService,
	orderService *services.OrderService,
) PaymentHandler {
	return PaymentHandler{
		service:      service,
		orderService: orderService,
	}
}

func (h *PaymentHandler) RegisterRoutes(r *gin.Engine) {
	orders := r.Group("/orders")
	{
		orders.POST("/:id/payments", h.Create)
		orders.GET("/:id/payments", h.GetByOrderID)
	}

	payments := r.Group("/payments")
	{
		payments.GET("/:id", h.GetByID)
	}
}

func (h *PaymentHandler) Create(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id is not a number"})
		return
	}

	var req models.PaymentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.GetOrder(uint(orderID))
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "заказ не найден" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	payment, updatedOrder, err := h.service.Create(uint(orderID), order.UserID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"payment": payment,
		"order":   updatedOrder,
	})
}

func (h *PaymentHandler) GetByOrderID(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id is not a number"})
		return
	}

	order, err := h.orderService.GetOrder(uint(orderID))
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "заказ не найден" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	payments, err := h.service.GetByOrderID(uint(orderID), order.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

func (h *PaymentHandler) GetByID(c *gin.Context) {
	paymentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment id is not a number"})
		return
	}

	payment, err := h.service.GetByID(uint(paymentID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}
