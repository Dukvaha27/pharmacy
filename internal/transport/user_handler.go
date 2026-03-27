package transport

import (
	"net/http"
	"pharmacy/internal/models"
	"pharmacy/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(userService services.UserService) UserHandler {
	return UserHandler{service: userService}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.POST("")
		users.PATCH("id")
		users.GET("/:id")
	}
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "UserID must be a number",
		})
		return
	}
	user, err := h.service.GetByID(uint64(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Update(c *gin.Context) {
	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Must be a normal data",
			"details": err.Error(),
		})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "UserID must be a number",
		})
		return
	}
	if err :=h.service.Update(uint64(id),&req); err!=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Updating user error",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated succesfully",
		"data": req,
	})

}

func (h *UserHandler) Create(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}
	if err := h.service.Create(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created succesfully",
		"data":    req,
	})
}
