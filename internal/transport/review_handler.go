package transport

import (
	"net/http"
	"pharmacy/internal/models"
	"pharmacy/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	service services.ReviewService
}

func NewReviewHandler(service services.ReviewService) ReviewHandler {
	return ReviewHandler{
		service: service,
	}
}

func (h *ReviewHandler) RegisterRoutes(r *gin.Engine) {
	reviews := r.Group("/reviews")
	{
		reviews.PATCH("/:id", h.Update)
		reviews.DELETE("/:id", h.Delete)
	}

	medicines := r.Group("/medicines")
	{
		medicines.GET("/:id/reviews", h.GetAll)
		medicines.POST("/:id/reviews", h.Create)
	}
}

func (h *ReviewHandler) GetAll(ctx *gin.Context) {
	medicineID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "medicine id is not a number"})
		return
	}

	reviews, err := h.service.GetAll(medicineID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reviews)
}

func (h *ReviewHandler) Create(ctx *gin.Context) {
	medicineID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "medicine id is not a number"})
		return
	}

	var req models.ReviewCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.MedicineID = uint(medicineID)

	if err := h.service.Create(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "review created",
	})
}

func (h *ReviewHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "review id is not a number"})
		return
	}

	var req models.ReviewUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(id, req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "review updated",
	})
}

func (h *ReviewHandler) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "review id is not a number"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "review deleted",
	})
}
