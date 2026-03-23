package transport

import (
	"net/http"
	"pharmacy/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SubCategoryHandler struct {
	services services.SubCategoryService
}

func NewSubCategoryHandler(subCategoryService services.SubCategoryService) SubCategoryHandler {
	return SubCategoryHandler{
		services: subCategoryService,
	}
}

func (h *SubCategoryHandler) RegisterRoutes(r *gin.Engine) {
	subCategories := r.Group("/subCategories")
	{
		subCategories.GET("", h.GetAllSubCategories)
		subCategories.GET("/with-parents/:id", h.GetCategoryWithSubCategories)
		subCategories.GET("/by-parentId/:categoryId", h.GetSubCategoriesByCategory)
		subCategories.POST("", h.CreateSubCategory)
		subCategories.GET("/:id", h.GetSubCategory)
		subCategories.PATCH("/:id", h.UpdateSubCategory)
		subCategories.DELETE("/:id", h.DeleteSubCategory)
	}
}

func (h *SubCategoryHandler) CreateSubCategory(c *gin.Context) {
	var req struct {
		CategoryID uint   `json:"category_id" binding:"required,min=1"`
		Name       string `json:"name" binding:"required,min=1,max=255"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}

	subCategory, err := h.services.Create(req.CategoryID, req.Name)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "category not found" {
			status = http.StatusNotFound
		} else if err.Error() == "subcategory with this name already exists in the category" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "subcategory created successfully",
		"data":    subCategory,
	})
}

func (h *SubCategoryHandler) GetSubCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subcategory id"})
		return
	}

	if id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subCategory id"})
		return
	}

	subCategory, err := h.services.GetByID(uint(id))
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "invalid subcategory id" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subCategory})
}

func (h *SubCategoryHandler) GetAllSubCategories(c *gin.Context) {
	subCategories, err := h.services.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  subCategories,
		"count": len(subCategories),
	})
}

func (h *SubCategoryHandler) GetSubCategoriesByCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("categoryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	if categoryID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	subCategories, err := h.services.GetByCategoryID(uint(categoryID))
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "invalid category id" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  subCategories,
		"count": len(subCategories),
	})
}

func (h *SubCategoryHandler) UpdateSubCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subcategory id"})
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

	subCategory, err := h.services.Update(uint(id), req.Name)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "subcategory not found" {
			status = http.StatusNotFound
		} else if err.Error() == "subcategory with this name already exists in the category" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "subcategory updated successfully",
		"data":    subCategory,
	})
}

func (h *SubCategoryHandler) DeleteSubCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subcategory id"})
		return
	}

	if id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subCategory id"})
		return
	}

	if err := h.services.Delete(uint(id)); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "subcategory not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "subcategory deleted successfully",
	})
}

func (h *SubCategoryHandler) GetCategoryWithSubCategories(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	if id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subCategory id"})
		return
	}

	category, subCategories, err := h.services.GetCategoryWithSubCategories(uint(id))
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "invalid category id" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category":      category,
		"subcategories": subCategories,
		"count":         len(subCategories),
	})
}
