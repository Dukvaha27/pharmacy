package main

import (
	"log"
	"pharmacy/internal/config"
	"pharmacy/internal/models"
	"pharmacy/internal/repository"
	"pharmacy/internal/services"
	"pharmacy/internal/transport"
	"pharmacy/internal/transport/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.SetUpDatabaseConnection()

	if err := db.AutoMigrate(
		&models.Cart{},
		&models.CartItem{},
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
	reviewRepo := repository.NewReviewRepository(db)
	cartRepo := repository.NewCartRepository(db)
	userRepo := repository.NewUserRepository(db)

	categoryService := services.NewCategoryService(categoryRepo)
	subCategoryService := services.NewSubCategoryService(subCategoryRepo, categoryRepo)
	medicineService := services.NewMedicineService(medicineRepo, categoryRepo, subCategoryRepo)
	reviewService := services.NewReviewService(reviewRepo, medicineRepo)
	cartService := services.NewCartService(cartRepo, userRepo, medicineRepo)
	userService := services.NewUserService(userRepo)

	router := gin.Default()

	limiter := middlewares.NewRateLimiter(config.RateLimitRPS, config.RateLimitBurst)
	router.Use(limiter.RateLimitMiddleware())

	transport.RegisterRoutes(router, categoryService, subCategoryService, medicineService, cartService, userService, reviewService)

	if err := router.Run(); err != nil {
		log.Fatalf("не удалось запустить HTTP-сервер: %v", err)
	}
}
