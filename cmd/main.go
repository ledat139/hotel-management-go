package main

import (
	"hotel-management/database"
	"hotel-management/router"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	database.InitDB()
	database.AutoMigrate()

	r := gin.Default()
	router.SetupRoutes(r)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
