package main

import (
	"app/api"
	"app/config"
	"app/pkg/logs"
	"app/pkg/smtp"
	"app/pkg/start"
	"app/storage/filestore"
	"app/storage/postgresql"
	redis_service "app/storage/redis"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(".env file not found" + err.Error())
	}

	cfg := config.Load()
	var loggerLevel string

	switch cfg.Env {
	case config.LocalMode:
		loggerLevel = logs.LevelDebug
	default:
		loggerLevel = logs.LevelInfo
	}

	log := logs.NewLogger("app", loggerLevel)
	defer func() {
		if err := logs.Cleanup(log); err != nil {
			return
		}
	}()

	store, err := postgresql.NewPostgresStore(cfg.DB, log)
	if err != nil {
		log.Panic("Error connect to postgres", logs.Error(err))
		return
	}
	defer store.Close()

	smtpService := smtp.NewSMTPService(log, &cfg.SMTP)
	cacheService := redis_service.NewRedisStore(&cfg.Redis, log)
	fileService := filestore.NewFilestore(&cfg.FileStorage, log)

	if err := start.Init(&cfg.DB, log, false, store.Role(), store.User()); err != nil {
		log.Panic("could not run start init", logs.Error(err))
		return
	}

	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	api.NewApi(r, &cfg, store, log, smtpService, cacheService, fileService)

	go func() {
		if err := r.Run(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
			log.Panic("Error listening server", logs.Error(err))
			os.Exit(1)
		}
	}()

	tickerPing := time.NewTicker(2 * time.Minute)

	log.Debug("setting ticker for ping")
	go func() {
		for range tickerPing.C {
			sendRequest(fmt.Sprintf("%s/%s", cfg.Server.Public, "api/ping"))
		}
	}()

	log.Info("Server running on port", logs.String("addr", cfg.Server.Host+cfg.Server.Port))
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	tickerPing.Stop()

	log.Info("db closed")
	store.Close()
}

func sendRequest(url string) {
	// Send the GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	} else {
		defer resp.Body.Close()
	}
}
