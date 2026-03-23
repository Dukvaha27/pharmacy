package transport

import (
	"net/http"
	"pharmacy/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) CategoryHandler {
	return CategoryHandler{
		service: categoryService,
	}
}

func (h *CategoryHandler) RegisterRoutes(r *gin.Engine) {
	categories := r.Group("/categories")
	{
		categories.GET("", h.GetAllCategories)
		categories.POST("", h.CreateCategory)
		categories.GET("/:id", h.GetCategory)
		categories.PATCH("/:id", h.UpdateCategory)
		categories.DELETE("/:id", h.DeleteCategory)
	}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required,min=1,max=255"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}

	category, err := h.service.Create(req.Name)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "category with this name already exists" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "category created successfully",
		"data":    category,
	})
}

func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	if id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	category, err := h.service.GetByID(uint(id))
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "invalid category id" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}

func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  categories,
		"count": len(categories),
	})
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required,min=1,max=255"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}

	if id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	category, err := h.service.Update(uint(id), req.Name)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "category not found" {
			status = http.StatusNotFound
		} else if err.Error() == "category with this name already exists" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "category updated successfully",
		"data":    category,
	})
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	if id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "category not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "category deleted successfully",
	})
}
