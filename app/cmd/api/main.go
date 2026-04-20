package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hotel.com/app/internal/config"
	"hotel.com/app/internal/database"
	"hotel.com/app/internal/handler"
	"hotel.com/app/internal/logging"
	"hotel.com/app/internal/repo"
	"hotel.com/app/internal/service"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Println("failed to load config:", err)
		os.Exit(1)
	}

	//create logger
	l := logging.New()
	l.Info("App initiated")

	//db connection
	db, err := database.NewConn(os.Getenv("USER_SERVICE_DATABASE_URL"))
	if err != nil {
		l.Error("Conection to database failed", "err", err)
		os.Exit(-1)
	}
	l.Info("Database connection successful")

	defer db.Close()

	err = database.RunMigrations(os.Getenv("USER_SERVICE_DATABASE_URL"), l)
	if err != nil {
		os.Exit(-1)
	}

	//jwt key file check
	if _, err := os.Stat("private.pem"); os.IsNotExist(err) {
		l.Error("JWT private key file not found", "err", err)
		os.Exit(-1)
	}
	//jwt key file check
	if _, err := os.Stat("public.pem"); os.IsNotExist(err) {
		l.Error("JWT private key file not found", "err", err)
		os.Exit(-1)
	}

	//repo creation
	r := repo.NewPostgreRepo(db)

	//service creation
	svc := service.New(l, r)

	// handler creation
	jwtConfig := handler.JWTConfig{
		Secret:     cfg.Server.Host, // This should be a secure secret in production
		Issuer:     "users-service",
		Expiration: 24 * time.Minute,
	}
	jwtAuth := handler.NewJWTAuthenticator(jwtConfig)
	h := handler.New(svc, l, jwtAuth)

	// server creation
	mux := h.NewServerMux(nil)
	port := cfg.Server.Port
	if port == 0 {
		port = 8080
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	l.Info("Server listening", "addr", srv.Addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Error("server failed", "err", err)
			os.Exit(1)
		}
	}()

	// Block until SIGTERM or SIGINT
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	l.Info("Shutting down server...")

	// Give in-flight requests 30s to finish
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		l.Error("Server forced to shutdown", "err", err)
	}

	l.Info("Server stopped")

}

// configuration is loaded from app/internal/config package
