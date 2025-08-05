package router

import (
	"hotel-management/database"
	"hotel-management/internal/handler"
	"hotel-management/internal/handler/admin"
	"hotel-management/internal/middleware"
	"hotel-management/internal/repository"
	"hotel-management/internal/usecase"
	"hotel-management/internal/usecase/admin_usecase"

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

	//Admin route
	statRepository := repository.NewStatRepository(database.DB)
	adminAuthUseCase := admin_usecase.NewAuthUseCase(userRepository)
	statUseCase := admin_usecase.NewStatUseCase(statRepository)
	adminHandler := admin.NewAdminHandler(adminAuthUseCase, statUseCase)

	roomRepository := repository.NewRoomRepository(database.DB)
	reviewRepository := repository.NewReviewRepository(database.DB)
	bookingRepository := repository.NewBookingRepository(database.DB)
	roomAdminUseCase := admin_usecase.NewRoomUseCase(roomRepository, bookingRepository, reviewRepository)
	roomAdminHandler := admin.NewRoomHandler(roomAdminUseCase)
	adminBookingUseCase := admin_usecase.NewBookingUseCase(bookingRepository)
	adminBookingHandler := admin.NewAdminBookingHandler(adminBookingUseCase)
	billRepository := repository.NewBillRepository(database.DB)
	billUseCase := admin_usecase.NewBillUseCase(billRepository)
	billHandler := admin.NewBillHandler(billUseCase)
	staffUseCase := admin_usecase.NewStaffUseCase(userRepository)
	staffHandler := admin.NewStaffHandler(staffUseCase)
	adminGroup := r.Group("/admin")
	{
		adminGroup.GET("/", middleware.RequireLogin(), middleware.RequireRoles("admin", "staff"), adminHandler.AdminDashboard)
		adminGroup.GET("/login", adminHandler.AdminLoginPage)
		adminGroup.POST("/login", adminHandler.HandleLogin)
		adminGroup.GET("/logout", adminHandler.HandleLogout)

		adminGroup.GET("/rooms", middleware.RequireRoles("admin", "staff"), roomAdminHandler.RoomManagementPage)
		adminGroup.GET("/rooms/create", middleware.RequireRoles("admin", "staff"), roomAdminHandler.CreateRoomPage)
		adminGroup.POST("/rooms/create", middleware.RequireRoles("admin", "staff"), roomAdminHandler.CreateRoom)
		adminGroup.GET("/rooms/:id", middleware.RequireRoles("admin", "staff"), roomAdminHandler.RoomDetailPage)
		adminGroup.GET("/rooms/edit/:id", middleware.RequireRoles("admin", "staff"), roomAdminHandler.EditRoomPage)
		adminGroup.POST("/rooms/edit/:id", middleware.RequireRoles("admin", "staff"), roomAdminHandler.UpdateRoom)
		adminGroup.POST("/rooms/delete/:id", middleware.RequireRoles("admin", "staff"), roomAdminHandler.DeleteRoom)
		adminGroup.GET("/bookings", middleware.RequireRoles("admin", "staff"), adminBookingHandler.ListBookings)
		adminGroup.GET("/bookings/:id", middleware.RequireRoles("admin", "staff"), adminBookingHandler.GetBookingDetail)
		adminGroup.GET("/bookings/edit/:id", middleware.RequireRoles("admin", "staff"), adminBookingHandler.EditBookingPage)
		adminGroup.POST("/bookings/edit/:id", middleware.RequireRoles("admin", "staff"), adminBookingHandler.EditBookingStatus)

		adminGroup.GET("/bills", middleware.RequireRoles("admin", "staff"), billHandler.ListBills)

		adminGroup.GET("/staffs", middleware.RequireRoles("admin"), staffHandler.ListStaffs)
		adminGroup.GET("/staffs/create", middleware.RequireRoles("admin"), staffHandler.CreateStaffPage)
		adminGroup.POST("/staffs/create", middleware.RequireRoles("admin"), staffHandler.CreateStaff)
		adminGroup.GET("/staffs/edit/:id", middleware.RequireRoles("admin"), staffHandler.EditStaffPage)
		adminGroup.POST("/staffs/edit/:id", middleware.RequireRoles("admin"), staffHandler.EditStaff)
		adminGroup.POST("/staffs/delete/:id", middleware.RequireRoles("admin"), staffHandler.DeleteStaff)

		adminGroup.GET("/customers", middleware.RequireRoles("admin"), staffHandler.ListCustomers)
	}
	//User routes
	userHandler := handler.NewUserHandler(userUseCase)
	r.PUT("/users/update-profile", middleware.RequireAuth(userRepository), userHandler.UpdateProfile)

	//Room routes
	roomUseCase := usecase.NewRoomUseCase(roomRepository)
	roomHandler := handler.NewRoomHandler(roomUseCase)
	r.POST("/rooms/search", roomHandler.FindAvailableRoom)

	//Booking routes
	paymentRepository := repository.NewPaymentRepository(database.DB)
	paymentUseCase := usecase.NewPaymentUseCase(paymentRepository, bookingRepository, billRepository)
	bookingUseCase := usecase.NewBookingUseCase(bookingRepository, *paymentUseCase)
	bookingHandler := handler.NewBookingHandler(bookingUseCase)
	bookingGroup := r.Group("/bookings")
	{
		bookingGroup.POST("/", middleware.RequireAuth(userRepository), bookingHandler.CreateBooking)
		bookingGroup.GET("/history", middleware.RequireAuth(userRepository), bookingHandler.GetBookingHistory)
		bookingGroup.GET("/:id/cancel", middleware.RequireAuth(userRepository), bookingHandler.CancelBooking)
	}
	//Review
	reviewUseCase := usecase.NewReviewUseCase(bookingRepository, reviewRepository)
	reviewHandler := handler.NewReviewHandler(reviewUseCase)
	r.POST("/reviews", middleware.RequireAuth(userRepository), reviewHandler.CreateReview)

	//Payment routes

	paymentHandler := handler.NewPaymentHandler(paymentUseCase)
	paymentGroup := r.Group("/payments")
	{
		paymentGroup.GET("/vnpay_return", paymentHandler.HandleVnpayCallback)
	}
}
