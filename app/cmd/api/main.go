package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"hotel.com/app/internal/database"
	"hotel.com/app/internal/handler"
	"hotel.com/app/internal/helper"
	"hotel.com/app/internal/logging"
	"hotel.com/app/internal/repo"
	"hotel.com/app/internal/service"
)

type Config struct {
	Port               string
	DbUrl              string
	LogLevel           string
	JWTSecret          string
	JWTIssuer          string
	JWTExpirationHours string
}

func main() {
	config := loadConfig()

	//create logger
	l := logging.New()
	l.Info("App initiated")

	//db connection
	db, err := database.NewConn(config.DbUrl)
	if err != nil {
		l.Error("Conection to database failed", "err", err)
		os.Exit(-1)
	}
	l.Info("Database connection successful")

	defer db.Close()

	err = database.RunMigrations(config.DbUrl, l)
	if err != nil {
		os.Exit(-1)
	}

	// JWT authenticator
	if config.JWTSecret == "" {
		l.Error("JWT secret not configured")
		os.Exit(-1)
	}

	expHrs, err := strconv.Atoi(config.JWTExpirationHours)
	if err != nil {
		l.Error("JWT expiration parse error, defaulting to 24h", "err", err.Error())
		expHrs = 24
	}

	jwtConfig := handler.JWTConfig{
		Secret:     config.JWTSecret,
		Issuer:     config.JWTIssuer,
		Expiration: time.Duration(expHrs) * time.Hour,
	}
	jwtAuth := handler.NewJWTAuthenticator(jwtConfig)

	//repo creation
	r := repo.NewPostgreRepo(db)

	//service creation
	svc := service.New(l, r)

	// handler creation
	h := handler.New(svc, l, jwtAuth)

	// server creation
	mux := h.NewServerMux(jwtAuth, nil)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
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

func loadConfig() *Config {
	return &Config{
		Port:               helper.Getenv("APP_PORT", "8080"),
		LogLevel:           helper.Getenv("LOG_LEVEL", "0"),
		DbUrl:              os.Getenv("DB_URL"),
		JWTSecret:          helper.Getenv("JWT_SECRET", ""),
		JWTIssuer:          helper.Getenv("JWT_ISSUER", ""),
		JWTExpirationHours: helper.Getenv("JWT_EXPIRATION_HOURS", "24"),
	}
}
