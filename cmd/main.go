package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/internal/config"
	"url-shortener/internal/db"
	"url-shortener/internal/server"
	"url-shortener/internal/shorten"
	"url-shortener/internal/storage/shortening"
	"url-shortener/internal/model"
)

func main() {
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer dbCancel()

	dbClient, err := db.Connect(dbCtx, config.LoadConfig().DB.DSN)
	if err != nil {
		log.Fatal("Failed to connect database :", err)
	}

	defer func() {
		if err := dbClient.Close(dbCtx); err != nil {
			log.Println("Failed to close database connection:", err)
		}
	}()

	if err := dbClient.Client().AutoMigrate(&model.Shortening{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	shorteningStorage := shortening.NewGormDB(dbClient.Client())
	shortener := shorten.NewService(shorteningStorage)

	srv := server.New(shortener)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(config.LoadConfig().ListenAddr(), srv); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Error running server: ", err)
		}
	}()

	log.Println("server started")
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Error closing server: ", err)
	}
	log.Println("Server stopped")
}
