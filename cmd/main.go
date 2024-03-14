package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/citadel-corp/shopifyx-marketplace/internal/common/db"
	"github.com/citadel-corp/shopifyx-marketplace/internal/common/middleware"
	"github.com/citadel-corp/shopifyx-marketplace/internal/product"
	"github.com/citadel-corp/shopifyx-marketplace/internal/user"
)

func main() {
	slogHandler := slog.NewTextHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(slogHandler))

	// Connect to database
	db, err := db.Connect(os.Getenv("DB_URL"))
	if err != nil {
		slog.Error(fmt.Sprintf("Cannot connect to database: %v", err))
	}

	// Create migrations
	err = db.UpMigration()
	if err != nil {
		slog.Error(fmt.Sprintf("Up migration failed: %v", err))
	}

	// initialize user domain
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	// initialize product domain
	productRepository := product.NewRepository(db)
	productService := product.NewService(productRepository)
	productHandler := product.NewHandler(productService)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text")
		io.WriteString(w, "Service ready")
	})

	mux.HandleFunc("POST /v1/user/register", middleware.PanicRecoverer(userHandler.CreateUser))
	mux.HandleFunc("POST /v1/user/login", middleware.PanicRecoverer(userHandler.Login))
	mux.HandleFunc("POST /v1/product", middleware.Authenticate(productHandler.CreateProduct))
	mux.HandleFunc("GET /v1/product", middleware.Authenticate(productHandler.GetProductList))

	httpServer := &http.Server{
		Addr:     ":8000",
		Handler:  mux,
		ErrorLog: slog.NewLogLogger(slogHandler, slog.LevelError),
	}

	go func() {
		slog.Info(fmt.Sprintf("HTTP server listening on %s", httpServer.Addr))
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error(fmt.Sprintf("HTTP server error: %v", err))
		}
		slog.Info("Stopped serving new connections.")
	}()

	// Listen for the termination signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Block until termination signal received
	<-stop
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	slog.Info(fmt.Sprintf("Shutting down HTTP server listening on %s", httpServer.Addr))
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP server shutdown error: %v", err)
	}
	slog.Info("Shutdown complete.")
}
