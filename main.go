package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dilaragorum/ticket-api/internal/ticket/database"

	_ "github.com/dilaragorum/ticket-api/docs"
	"github.com/dilaragorum/ticket-api/internal/ticket/handler"
	"github.com/dilaragorum/ticket-api/internal/ticket/repository"
	"github.com/dilaragorum/ticket-api/internal/ticket/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Ticket API
// @version 1.0
// @description TicketService

// @contact.name Dilara Görüm
// @contact.email dilaragorum@gmail.com

// @host localhost:3000
func main() {
	e := echo.New()

	err := godotenv.Load(".env.dev")
	if err != nil {
		log.Fatal(err)
	}

	connectionPool, err := database.Setup()
	if err != nil {
		log.Fatal(err)
	}
	database.Migrate()

	ticketRepo := repository.NewDefaultRepository(connectionPool)
	ticketSvc := service.NewDefaultService(ticketRepo)
	handler.NewDefaultTicketHandler(e, ticketSvc)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	go func() {
		if err := e.Start(":3000"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //nolint:gomnd
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
