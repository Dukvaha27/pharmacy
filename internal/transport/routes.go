package transport

import (
	"pharmacy/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	categoryService services.CategoryService,
	subCategoryService services.SubCategoryService,
) {
	categoryHandler := NewCategoryHandler(categoryService)
	subCategoryHandler := NewSubCategoryHandler(subCategoryService)

	categoryHandler.RegisterRoutes(router)
	subCategoryHandler.RegisterRoutes(router)
}
