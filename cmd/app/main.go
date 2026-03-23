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
	); err != nil {
		log.Fatalf("не удалось выполнить миграции: %v", err)
	}

	categoryRepo := repository.NewCategoryRepository(db)
	subCategoryRepo := repository.NewSubCategoryRepository(db)

	categoryService := services.NewCategoryService(categoryRepo)
	subCategoryService := services.NewSubCategoryService(subCategoryRepo, categoryRepo)

	router := gin.Default()

	transport.RegisterRoutes(router, categoryService, subCategoryService)

	if err := router.Run(); err != nil {
		log.Fatalf("не удалось запустить HTTP-сервер: %v", err)
	}
}
