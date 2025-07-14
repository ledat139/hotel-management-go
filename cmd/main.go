package main

import (
	"log"
	"hotel-management/database"
	"hotel-management/router"
	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

	r := gin.Default()
	router.SetupRoutes(r)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
