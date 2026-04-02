package transport

import (
	"pharmacy/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	categoryService services.CategoryService,
	subCategoryService services.SubCategoryService,
	medicineService services.MedicineService,
	cartService services.CartService,
	userService services.UserService,
	reviewService services.ReviewService,
) {
	categoryHandler := NewCategoryHandler(categoryService)
	subCategoryHandler := NewSubCategoryHandler(subCategoryService)
	medicineHandler := NewMedicineHandler(medicineService)
	userHandler := NewUserHandler(userService)
	cartHandler := NewCartHandler(cartService)
	reviewHandler := NewReviewHandler(reviewService)

	categoryHandler.RegisterRoutes(router)
	subCategoryHandler.RegisterRoutes(router)
	medicineHandler.RegisterRoutes(router)
	cartHandler.RegisterRoutes(router)
	reviewHandler.RegisterRoutes(router)
	userHandler.RegisterRoutes(router)
}
