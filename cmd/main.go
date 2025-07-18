package main

import (
	"hotel-management/database"
	_ "hotel-management/docs"
	"hotel-management/internal/middleware"
	"hotel-management/internal/utils"
	"hotel-management/router"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title Hotel Management API
// @version 1.0
// @description This is an API for hotel management.

// @host localhost:8080
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	utils.InitJWT()
	database.InitDB()
	database.AutoMigrate()
	utils.InitI18n()

	r := gin.Default()
	r.Use(middleware.I18nMiddleware())
	router.SetupRoutes(r)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
