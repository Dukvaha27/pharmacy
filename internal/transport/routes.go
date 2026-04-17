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
	orderService *services.OrderService,
	paymentService *services.PaymentService,
	promocodeService *services.PromocodeService,
) {
	categoryHandler := NewCategoryHandler(categoryService)
	subCategoryHandler := NewSubCategoryHandler(subCategoryService)
	medicineHandler := NewMedicineHandler(medicineService)
	userHandler := NewUserHandler(userService, orderService)
	cartHandler := NewCartHandler(cartService)
	reviewHandler := NewReviewHandler(reviewService)
	orderHandler := NewOrderHandler(orderService)
	paymentHandler := NewPaymentHandler(paymentService, orderService)
	promocodeHandler := NewPromocodeHandler(promocodeService)

	categoryHandler.RegisterRoutes(router)
	subCategoryHandler.RegisterRoutes(router)
	medicineHandler.RegisterRoutes(router)
	cartHandler.RegisterRoutes(router)
	reviewHandler.RegisterRoutes(router)
	userHandler.RegisterRoutes(router)
	orderHandler.RegisterRoutes(router)
	paymentHandler.RegisterRoutes(router)
	promocodeHandler.RegisterRoutes(router)
}
