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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	bankaccount "github.com/citadel-corp/shopifyx-marketplace/internal/bank_account"
	"github.com/citadel-corp/shopifyx-marketplace/internal/common/db"
	"github.com/citadel-corp/shopifyx-marketplace/internal/common/middleware"
	"github.com/citadel-corp/shopifyx-marketplace/internal/image"
	"github.com/citadel-corp/shopifyx-marketplace/internal/product"
	"github.com/citadel-corp/shopifyx-marketplace/internal/user"
	"github.com/gorilla/mux"
)

func main() {
	slogHandler := slog.NewTextHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(slogHandler))

	// Connect to database
	env := os.Getenv("ENV")
	sslMode := "disable"
	if env == "production" {
		sslMode = "verify-full sslrootcert=ap-southeast-1-bundle.pem"
	}
	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), sslMode)
	db, err := db.Connect(dbURL)
	if err != nil {
		slog.Error(fmt.Sprintf("Cannot connect to database: %v", err))
		os.Exit(1)
	}

	// Create migrations
	// err = db.UpMigration()
	// if err != nil {
	// 	slog.Error(fmt.Sprintf("Up migration failed: %v", err))
	// 	os.Exit(1)
	// }

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

	// initialize image domain
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("S3_ID"), os.Getenv("S3_SECRET_KEY"), ""),
	})
	if err != nil {
		slog.Error(fmt.Sprintf("Cannot create AWS session: %v", err))
		os.Exit(1)
	}
	imageService := image.NewService(sess)
	imageHandler := image.NewHandler(imageService)

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
	pr.HandleFunc("/{productId}", middleware.PanicRecoverer(middleware.Authorized(productHandler.DeleteProduct))).Methods(http.MethodDelete)
	pr.HandleFunc("/{productId}/buy", middleware.PanicRecoverer(middleware.Authorized(productHandler.PurchaseProduct))).Methods(http.MethodPost)
	pr.HandleFunc("/{productId}/stock", middleware.PanicRecoverer(middleware.Authorized(productHandler.UpdateStockProduct))).Methods(http.MethodPost)

	// bank routes
	br := v1.PathPrefix("/bank").Subrouter()
	br.HandleFunc("/account", middleware.PanicRecoverer(middleware.Authorized(bankAccountHandler.CreateBankAccount))).Methods(http.MethodPost)
	br.HandleFunc("/account", middleware.PanicRecoverer(middleware.Authorized(bankAccountHandler.ListBankAccount))).Methods(http.MethodGet)
	br.HandleFunc("/account", middleware.PanicRecoverer(middleware.Authorized(bankAccountHandler.PartialUpdateBankAccount))).Methods(http.MethodPatch)
	br.HandleFunc("/account/{uuid}", middleware.PanicRecoverer(middleware.Authorized(bankAccountHandler.PartialUpdateBankAccount))).Methods(http.MethodPatch)
	br.HandleFunc("/account/{uuid}", middleware.PanicRecoverer(middleware.Authorized(bankAccountHandler.DeleteBankAccount))).Methods(http.MethodDelete)

	// image routes
	ir := v1.PathPrefix("/image").Subrouter()
	ir.HandleFunc("", middleware.PanicRecoverer(middleware.Authorized(imageHandler.UploadToS3))).Methods(http.MethodPost)

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
