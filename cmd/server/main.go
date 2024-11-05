package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"GO-API/internal/infrastructure/database/postgres"
	"GO-API/internal/infrastructure/processor"
	"GO-API/internal/interface/handler"
	"GO-API/internal/interface/middleware"
	"GO-API/internal/pkg/logger"
	"GO-API/internal/usecase"
)

func main() {
	config := &postgres.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     5432,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "go_api"),
		SSLMode:  "disable",
	}

	//start transaction to connect DB
	db, err := postgres.NewConnection(config)
	if err != nil {
		logger.Error("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	logger.Info("Successfully connected to database")
	defer db.Close()

	paymentRepo := postgres.NewPaymentRepository(db)

	if err := paymentRepo.InitTable(); err != nil {
		log.Fatalf("Failed to init tables: %v", err)
	}

	paymentProcessor := processor.NewPaymentProcessor()

	paymentUseCase := usecase.NewPaymentUseCase(paymentRepo, paymentProcessor)

	paymentHandler := handler.NewPaymentHandler(paymentUseCase)

	router := mux.NewRouter()
	router.Use(middleware.CORS)
	router.Use(middleware.RequestLogger)

	router.HandleFunc("/health", handler.HealthCheck).Methods(http.MethodGet)

	paymentHandler.RegisterRoutes(router)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", getEnv("PORT", "8080")),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info("Server starting on port %s", srv.Addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error("Server failed to start: %v", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Info("Server forces to shutdown: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
