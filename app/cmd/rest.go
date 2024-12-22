package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"payment-gateway/internal/rest"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
)

const (
	ctxTimeoutSec      = 60
	shutdownTimeoutSec = 185
)

var restCommand = &cobra.Command{
	Use:   "rest",
	Short: "Start REST server",
	Run:   restServer,
}

func init() {
	rootCmd.AddCommand(restCommand)
}

func initTimeoutCtx() time.Duration {
	timeoutCtx := time.Second * time.Duration(ctxTimeoutSec)
	return timeoutCtx
}

func restServer(cmd *cobra.Command, args []string) {
	e := echo.New()

	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	timeoutCtx := initTimeoutCtx()

	//Run cron in the same process as web server
	InitCron()

	registerControllers(e, timeoutCtx)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	srvAddress := os.Getenv("SERVER_ADDRESS")

	go func() {
		if err := e.Start(srvAddress); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(fmt.Sprintf("shutting down the server: %s", err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 185 seconds.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Shutdown timeout (185 seconds)
	shutdownTimeout := 185 * time.Second

	<-stop
	log.Printf("Shutting down server...\n")

	// Create a shutdown context with timeout
	gracefulCtx, cancelShutdown := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelShutdown()

	// Attempt to gracefully shutdown the server
	if err := e.Shutdown(gracefulCtx); err != nil {
		log.Fatal(fmt.Errorf("failed to gracefully shutdown the server. detailed error:\n%s", err))
	} else {
		log.Printf("Server gracefully shutdown ...\n")
	}
}

func registerControllers(e *echo.Echo, timeoutCtx time.Duration) {
	rest.InstallTransactionController(e, TransactionService, timeoutCtx)
}
