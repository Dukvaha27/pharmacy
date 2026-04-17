package transport

import (
	"net/http"
	"pharmacy/internal/models"
	"pharmacy/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	service services.CartService
}

func NewCartHandler(cartService services.CartService) CartHandler {
	return CartHandler{service: cartService}
}

func (h *CartHandler) RegisterRoutes(r *gin.Engine) {
	usersCart := r.Group("/users/:id/cart")
	{
		usersCart.GET("", h.GetByUserID)
		usersCart.DELETE("", h.ClearCart)

		usersCart.POST("/items", h.AddItem)
		usersCart.PATCH("/items/:item_id", h.UpdateItem)
		usersCart.DELETE("/items/:item_id", h.DeleteItem)
	}
}

func (h *CartHandler) GetByUserID(c *gin.Context) {
	userID, err := parseUint64Param(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id is not a number"})
		return
	}

	cart, err := h.service.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) AddItem(c *gin.Context) {
	userID, err := parseUint64Param(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id is not a number"})
		return
	}

	var req models.CartItemCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddItem(userID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.service.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) UpdateItem(c *gin.Context) {
	userID, err := parseUint64Param(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id is not a number"})
		return
	}

	itemID, err := parseUint64Param(c, "item_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "item id is not a number"})
		return
	}

	var req models.CartItemUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateItem(userID, itemID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.service.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) DeleteItem(c *gin.Context) {
	userID, err := parseUint64Param(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id is not a number"})
		return
	}

	itemID, err := parseUint64Param(c, "item_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "item id is not a number"})
		return
	}

	if err := h.service.DeleteItem(userID, itemID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.service.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, err := parseUint64Param(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id is not a number"})
		return
	}

	if err := h.service.ClearCart(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.service.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func parseUint64Param(c *gin.Context, name string) (uint64, error) {
	value, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}
