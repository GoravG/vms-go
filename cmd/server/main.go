package main

import (
	"fmt"
	"net/http"
	"time"

	"vms_go/internal/config"
	"vms_go/internal/db"
	"vms_go/internal/handlers"
	"vms_go/internal/migrations"
	"vms_go/internal/redis"
	"vms_go/internal/utils"
	"vms_go/internal/ws"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.FromEnv()

	utils.InitLogger(cfg.LOG_FILE)
	utils.LogInfo("Starting Application")
	defer utils.CloseLogger()

	connection_string := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.DB_USER, cfg.DB_PASS, cfg.DB_HOST, cfg.DB_PORT, cfg.DB_NAME)
	database := db.Connect(connection_string)
	migrator := migrations.NewMigrator(database)
	migrator.RunMigrtations()
	defer database.Close()

	authHandler := &handlers.AuthHandler{
		DB:         database,
		HMACSecret: cfg.HMAC_SECRET_KEY,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("POST /login", authHandler.Login)
	// mux.HandleFunc("POST /logout", authHandler.Logout)
	userHandler := &handlers.UserHandler{DB: database}
	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("GET /checkin", userHandler.Checkin)

	hub := ws.NewHub()
	mux.Handle("/ws", &handlers.WSHandler{Hub: hub})

	// protected := http.HandlerFunc(authHandler.Protected)
	// mux.Handle("GET /protected", middleware.Auth(protected, cfg.HMACSecret))
	redisConnectionString := fmt.Sprintf("%s:%s", cfg.REDIS_HOST, cfg.REDIS_PORT)
	go redis.SubscribeToRedis(hub, redisConnectionString, cfg.REDIS_CHANNEL)

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	utils.LogInfof("server listening on %s", cfg.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		utils.LogErrorf("server failed: %v", err)
	}
}
