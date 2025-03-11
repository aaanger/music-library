package main

import (
	"context"
	_ "github.com/aaanger/music-library/docs"
	"github.com/aaanger/music-library/internal/handler"
	"github.com/aaanger/music-library/internal/repository"
	"github.com/aaanger/music-library/internal/service"
	"github.com/aaanger/music-library/pkg/db"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Онлайн библиотека песен
// @version 1.0
// @description Swagger API для бибилотеки песен
// @BasePath /api/v1

type server struct {
	httpServer *http.Server
}

func (srv *server) run(port string, handler http.Handler) error {
	srv.httpServer = &http.Server{
		Addr:         port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return srv.httpServer.ListenAndServe()
}

func (srv *server) shutdown(ctx context.Context) error {
	return srv.httpServer.Shutdown(ctx)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("Error loading .env file: %s", err)
	}

	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = logrus.InfoLevel
	}

	log := logrus.New()
	log.SetLevel(logLevel)

	db, err := db.Open(db.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		Username: os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		DBName:   os.Getenv("PSQL_DBNAME"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	})
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}

	apiURL := os.Getenv("API_URL")

	repo := repository.NewMusicRepository(db, log)
	service := service.NewMusicService(repo, apiURL, log)
	handler := handler.NewMusicHandler(service, log)

	srv := new(server)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		err := srv.run(":"+port, handler.InitRoutes())
		if err != nil {
			log.Fatalf("Error running the server: %s", err)
		}
		log.Infof("Server running on port :%s", port)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	err = srv.shutdown(context.Background())
	log.Infof("Shutting down the server")
	if err != nil {
		log.Errorf("Error shutting down the server: %s", err)
	}

	err = db.Close()
	if err != nil {
		log.Errorf("Error closing database: %s", err)
	}
}
