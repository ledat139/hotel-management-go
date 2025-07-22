package main

import (
	"hotel-management/database"
	_ "hotel-management/docs"
	"hotel-management/internal/middleware"
	"hotel-management/internal/utils"
	"hotel-management/router"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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
	utils.InitMail()
	utils.InitGoogleAuth()

	r := gin.Default()
	r.Use(middleware.I18nMiddleware())
	//Load static and template
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "..")

	assetPath := filepath.Join(basePath, "web", "assets")
	r.Static("/assets", assetPath)

	r.LoadHTMLGlob(filepath.Join(basePath, "/web/templates/**/*.html"))

	store := cookie.NewStore([]byte(os.Getenv("SECRET_KEY")))
	r.Use(sessions.Sessions("mysession", store))
	router.SetupRoutes(r)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
