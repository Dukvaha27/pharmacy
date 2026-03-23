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
) {
	categoryHandler := NewCategoryHandler(categoryService)
	subCategoryHandler := NewSubCategoryHandler(subCategoryService)
	medicineHandler := NewMedicineHandler(medicineService)

	categoryHandler.RegisterRoutes(router)
	subCategoryHandler.RegisterRoutes(router)
	medicineHandler.RegisterRoutes(router)
}
