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

	"github.com/Vaibhavsahu2810/redis-go/api"
	"github.com/Vaibhavsahu2810/redis-go/internal/config"
	"github.com/Vaibhavsahu2810/redis-go/internal/email"
	"github.com/Vaibhavsahu2810/redis-go/internal/queue"
	"github.com/Vaibhavsahu2810/redis-go/internal/templates"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	cfg := config.New()

	tmpl, err := templates.New()
	if err != nil {
		log.Fatalf("Error initializing templates: %v", err)
	}

	redisClient, err := queue.NewRedisClient(cfg)
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	defer redisClient.Close()

	emailService := email.NewSender(cfg, tmpl)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go queue.StartWorker(ctx, redisClient, emailService)

	router := gin.Default()
	api.RegisterHandlers(router, redisClient)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Println("Server shut down successfully")
}
