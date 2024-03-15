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

	bankaccount "github.com/citadel-corp/shopifyx-marketplace/internal/bank_account"
	"github.com/citadel-corp/shopifyx-marketplace/internal/common/db"
	"github.com/citadel-corp/shopifyx-marketplace/internal/common/middleware"
	"github.com/citadel-corp/shopifyx-marketplace/internal/product"
	"github.com/citadel-corp/shopifyx-marketplace/internal/user"
	"github.com/gorilla/mux"
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

	// initialize bank account domain
	bankAccountRepository := bankaccount.NewRepository(db)
	bankAccountService := bankaccount.NewService(bankAccountRepository)
	bankAccountHandler := bankaccount.NewHandler(bankAccountService)

	// initialize product domain
	productRepository := product.NewRepository(db)
	productService := product.NewService(productRepository, userRepository, bankAccountRepository)
	productHandler := product.NewHandler(productService)

	r := mux.NewRouter()
	v1 := r.PathPrefix("/v1").Subrouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text")
		io.WriteString(w, "Service ready")
	})

	// user routes
	ur := v1.PathPrefix("/user").Subrouter()
	ur.HandleFunc("/register", middleware.PanicRecoverer(userHandler.CreateUser)).Methods(http.MethodPost)
	ur.HandleFunc("/login", middleware.PanicRecoverer(userHandler.Login)).Methods(http.MethodPost)

	// product routes
	pr := v1.PathPrefix("/product").Subrouter()
	pr.HandleFunc("", middleware.PanicRecoverer(middleware.Authorized(productHandler.CreateProduct))).Methods(http.MethodPost)
	pr.HandleFunc("", middleware.PanicRecoverer(middleware.Authenticate(productHandler.GetProductList))).Methods(http.MethodGet)
	pr.HandleFunc("/{productId}", middleware.PanicRecoverer(middleware.Authorized(productHandler.PatchProduct))).Methods(http.MethodPatch)
	pr.HandleFunc("/{productId}", middleware.PanicRecoverer(productHandler.GetProduct)).Methods(http.MethodGet)
	pr.HandleFunc("/{productId}/buy", middleware.PanicRecoverer(middleware.Authorized(productHandler.PurchaseProduct))).Methods(http.MethodPost)

	// bank routes
	br := v1.PathPrefix("/bank").Subrouter()
	br.HandleFunc("/account", middleware.PanicRecoverer(middleware.Authorized(bankAccountHandler.CreateBankAccount))).Methods(http.MethodPost)
	br.HandleFunc("/account", middleware.PanicRecoverer(middleware.Authorized(bankAccountHandler.ListBankAccount))).Methods(http.MethodGet)
	br.HandleFunc("/account/{uuid}", middleware.PanicRecoverer(middleware.Authorized(bankAccountHandler.PartialUpdateBankAccount))).Methods(http.MethodPatch)
	br.HandleFunc("/account/{uuid}", middleware.PanicRecoverer(middleware.Authorized(bankAccountHandler.DeleteBankAccount))).Methods(http.MethodDelete)

	httpServer := &http.Server{
		Addr:     ":8000",
		Handler:  r,
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
