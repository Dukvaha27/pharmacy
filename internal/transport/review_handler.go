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

func NewReviewHandler(reviewService services.ReviewService) ReviewHandler {
	return ReviewHandler{service: reviewService}
}

// GET /medicines/:id/reviews — список отзывов к лекарству.
// POST /medicines/:id/reviews — добавить отзыв.
// PATCH /reviews/:id — частично изменить отзыв (например, текст или оценку).
// DELETE /reviews/:id — удалить отзыв.

func (h *ReviewHandler) RegisterRoutes(r *gin.Engine) {
	medicinesByIDReview := r.Group("/medicines/:id/review")
	{
		medicinesByIDReview.GET("", h.GetAll)
		medicinesByIDReview.POST("", h.Create)
	}
	reviewsByID := r.Group("/reviews/:id")
	{
		reviewsByID.PATCH("", h.Update)
		reviewsByID.DELETE("", h.Delete)
	}
}

func (h *ReviewHandler) GetAll(c *gin.Context) {
	medicineID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Medicine ID is not a number",
			"details": err.Error(),
		})
		return
	}
	reviews, err := h.service.GetAll(uint64(medicineID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Ошибка при получении данных",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, reviews)
}

func (h *ReviewHandler) Create(c *gin.Context) {
	var req models.ReviewCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неправильно переданы данные",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.Create(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Ошибка при создании данных",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, "Review created")
}

func (h *ReviewHandler) Update(c *gin.Context) {
	var req models.ReviewUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неправильно переданы данные",
			"details": err.Error(),
		})
		return
	}

	reviewID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Review ID is not a number",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.Update(uint64(reviewID), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Ошибка при обновлении данных",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, "Review updated")
}

func (h *ReviewHandler) Delete(c *gin.Context) {
	reviewID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Review ID is not a number",
			"details": err.Error(),
		})
		return
	}

	if err := h.service.Delete(uint64(reviewID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Ошибка удаления отзыва",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, "Review deleted")
}
