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
		usersCart.POST("/items", h.AddItem)
		usersCart.PATCH("/items/:item_id", h.UpdateItem)
		usersCart.DELETE("/items/:item_id", h.Delete)
	}
}

func (h *CartHandler) Delete(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "User ID is not a number",
			"details": err.Error(),
		})
		return
	}
	itemID, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "item ID is not a number",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.DeleteItem(uint(userID), uint(itemID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ошибка при удалении позиции из твоей корзины.",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, "Item from your cart deleted")
}

func (h *CartHandler) UpdateItem(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "User ID is not a number",
			"details": err.Error(),
		})
		return
	}
	itemID, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "item ID is not a number",
			"details": err.Error(),
		})
		return
	}

	var req models.CartItemUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := h.service.UpdateItem(uint64(userID), uint64(itemID), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Ошибка обновления",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *CartHandler) GetByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "User ID is not a number",
			"details": err.Error(),
		})
		return
	}
	cart, err := h.service.GetByUserID(uint64(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Ошибка при взятии корзины пользователя.",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, cart)

}

func (h *CartHandler) AddItem(c *gin.Context) {
	var req models.CartItemCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "User ID is not a number",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.AddItem(uint(userID), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Ошибка при обработке данных",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Added item in your cart",
		"data":    req,
	})

}
