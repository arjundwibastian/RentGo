package main

import (
	"context"
	"log"
	"milestone-02/internal/config"
	"milestone-02/internal/domain"
	"milestone-02/internal/handler"
	"milestone-02/internal/infrastructure/database"
	"milestone-02/internal/middleware"
	"milestone-02/internal/repository/db"
	"milestone-02/internal/repository/http"
	"milestone-02/internal/usecase"
	"milestone-02/internal/validation"
	stdhttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	err := godotenv.Load("file.env")
	if err != nil {
		log.Fatal("Error loading file.env:", err)
	}
	cfg, err := config.Load() 
	if err != nil {
		log.Fatal("Error creating config:", err)
	}

	dbConn, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	jwtSecret := []byte(cfg.JWT.Secret)

	jwtMiddleware := middleware.NewJWTMiddleware(jwtSecret)

	repo := db.NewUserRepo(dbConn)
	emailService := http.NewSendGridEmailService(cfg.App.APIKey)

	authUC := usecase.NewAuthUseCase(repo,cfg.App.AdminKey, jwtSecret, cfg.JWT.Expiry)
	userUC := usecase.NewUserUseCase(repo, emailService)
	adminUC := usecase.NewAdminUsecase(repo)

	authHandler := handler.NewAuthHandler(authUC)
	userHandler := handler.NewUserHandler(userUC)
	adminHandler := handler.NewAdminHandler(adminUC)

	e := echo.New()
	e.Use(echoMiddleware.RequestLogger())
	e.Use(echoMiddleware.Recover())
	e.Validator = validation.New()

	e.POST("/api/v1/users/register", authHandler.RegisterHandler)
	e.POST("/api/v1/users/login", authHandler.LoginHandler)
	e.GET("/api/v1/swagger/*", echoSwagger.WrapHandler)


	user := e.Group("")
	user.Use(jwtMiddleware.Authenticate, jwtMiddleware.UserOnly)
	user.GET("/api/v1/users/profile", userHandler.GetProfile)
	user.POST("/api/v1/users/topup", userHandler.AddUserBalanceHandler)
	user.GET("/api/v1/vehicles", userHandler.GetVehicleListHandler)
	user.POST("/api/v1/vehicles/date", userHandler.GetVehicleAvailableByDateHandler)
	user.GET("/api/v1/bookings/history", userHandler.GetUserBookingHistoryHandler)
	user.POST("/api/v1/bookings", userHandler.CreateBookingHandler)
	user.POST("/api/v1/bookings/cancel", userHandler.CancelBookingHandler)


	admin := e.Group("")
	admin.Use(jwtMiddleware.Authenticate, jwtMiddleware.AdminOnly)
	admin.POST("/api/v1/vehicles", adminHandler.CreateNewVehicleHandler)
	admin.PUT("/api/v1/vehicles/:id", adminHandler.UpdateVehicleHandler)
	admin.GET("/api/v1/reports/revenue", adminHandler.GetRevenueReportHandler)
	admin.GET("/api/v1/reports/top-vehicle", adminHandler.GetTopVehicleHandler)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go startBookingCompletionWorker(ctx, repo)

	go func() {
		if err := e.Start(":" + cfg.App.Port); err != nil && err != stdhttp.ErrServerClosed {
			e.Logger.Fatal("shutting down the server: ", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		e.Logger.Fatal(err)
	}
}

func startBookingCompletionWorker(ctx context.Context, repo domain.UserRepository) {
	complete := func() {
		if err := repo.CompletePastBookings(); err != nil {
			log.Printf("complete past bookings failed: %v", err)
		}
	}

	complete()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			complete()
		case <-ctx.Done():
			return
		}
	}
}