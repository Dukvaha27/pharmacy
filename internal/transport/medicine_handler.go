package transport

import (
	"errors"
	"net/http"
	"pharmacy/internal/models"
	"pharmacy/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MedicineHandler struct {
	service services.MedicineService
}

func NewMedicineHandler(service services.MedicineService) *MedicineHandler {
	return &MedicineHandler{service: service}
}

func (h *MedicineHandler) RegisterRoutes(r *gin.Engine) {
	medicines := r.Group("/medicines")
	{
		medicines.GET("", h.List)
		medicines.POST("", h.Create)
		medicines.GET("/:id", h.Get)
		medicines.PATCH("/:id", h.Update)
		medicines.DELETE("/:id", h.Delete)
		medicines.POST("/:id/check-stock", h.CheckStock)
	}
}

func (h *MedicineHandler) Create(c *gin.Context) {
	var req models.MedicineCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	medicine, err := h.service.CreateMedicine(req)
	if err != nil {
		if errors.Is(err, services.ErrCategoryNotFound) ||
			errors.Is(err, services.ErrSubCategoryNotFound) ||
			errors.Is(err, services.ErrInvalidSubCategory) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, medicine)
}

func (h *MedicineHandler) Get(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор"})
		return
	}

	medicine, err := h.service.GetMedicineByID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrMedicineNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrInvalidID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, medicine)
}

func (h *MedicineHandler) Update(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор"})
		return
	}

	var req models.MedicineUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	medicine, err := h.service.UpdateMedicine(uint(id), req)
	if err != nil {
		if errors.Is(err, services.ErrMedicineNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrInvalidID) ||
			errors.Is(err, services.ErrCategoryNotFound) ||
			errors.Is(err, services.ErrSubCategoryNotFound) ||
			errors.Is(err, services.ErrInvalidSubCategory) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, medicine)
}

func (h *MedicineHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор"})
		return
	}

	if err := h.service.DeleteMedicine(uint(id)); err != nil {
		if errors.Is(err, services.ErrMedicineNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrInvalidID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *MedicineHandler) List(c *gin.Context) {
	var filter models.MedicineFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if filter.Limit > 100 {
		filter.Limit = 100
	}

	medicines, total, err := h.service.ListMedicines(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  medicines,
		"total": total,
		"page":  filter.Page,
		"limit": filter.Limit,
	})
}

func (h *MedicineHandler) CheckStock(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный идентификатор"})
		return
	}

	var req struct {
		Quantity int `json:"quantity" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.service.CheckStock(uint(id), req.Quantity); err != nil {
		if errors.Is(err, services.ErrMedicineNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrInsufficientStock) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "товар доступен в указанном количестве",
		"quantity": req.Quantity,
	})
}
