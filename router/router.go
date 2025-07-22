package router

import (
	"hotel-management/database"
	"hotel-management/internal/handler"
	"hotel-management/internal/repository"
	"hotel-management/internal/usecase"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth routes
	userRepository := repository.NewUserRepository(database.DB)
	userUseCase := usecase.NewUserUseCase(userRepository)
	authUseCase := usecase.NewAuthUseCase(userRepository)
	authHandler := handler.NewAuthHandler(userUseCase, authUseCase)
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh-token", authHandler.RefreshToken)
		authGroup.GET("/google/login", authHandler.GoogleLoginHandler)
		authGroup.GET("/google/callback", authHandler.GoogleCallbackHandler)
	}

	//Mail routes
	mailUseCase := usecase.NewMailUseCase(userRepository)
	mailHandler := handler.NewMailHandler(mailUseCase)
	mailGroup := r.Group("/mail")
	{
		mailGroup.POST("/smtp-verify", mailHandler.SendVerificationEmail)
		mailGroup.GET("/verify-account", mailHandler.ActiveAccountHandler)
		mailGroup.GET("/reset-password", mailHandler.ResetPassword)
	}

}
