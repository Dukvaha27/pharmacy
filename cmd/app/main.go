package main

import (
	"log"
	"pharmacy/internal/config"
	"pharmacy/internal/models"
	"pharmacy/internal/repository"
	"pharmacy/internal/services"
	"pharmacy/internal/transport"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.SetUpDatabaseConnection()

	if err := db.AutoMigrate(
		&models.CartItem{},
		&models.Cart{},
		&models.Category{},
		&models.Review{},
		&models.SubCategory{},
		&models.User{},
		&models.Medicine{},
	); err != nil {
		log.Fatalf("не удалось выполнить миграции: %v", err)
	}

	categoryRepo := repository.NewCategoryRepository(db)
	subCategoryRepo := repository.NewSubCategoryRepository(db)
	medicineRepo := repository.NewMedicineRepository(db)

	categoryService := services.NewCategoryService(categoryRepo)
	subCategoryService := services.NewSubCategoryService(subCategoryRepo, categoryRepo)
	medicineService := services.NewMedicineService(medicineRepo, categoryRepo, subCategoryRepo)

	router := gin.Default()

	transport.RegisterRoutes(router, categoryService, subCategoryService, medicineService)

	if err := router.Run(); err != nil {
		log.Fatalf("не удалось запустить HTTP-сервер: %v", err)
	}
}
