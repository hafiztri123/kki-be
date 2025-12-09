package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/hafiztri123/kki-be/internal/config"
	"github.com/hafiztri123/kki-be/internal/constants"
	"github.com/hafiztri123/kki-be/internal/handler"
	"github.com/hafiztri123/kki-be/internal/repository"
	"github.com/hafiztri123/kki-be/internal/service"
	"github.com/hafiztri123/kki-be/internal/utils"
	"github.com/joho/godotenv"
)

const (
	LogFileName = "logs.log"
	LogFilePath = "./logs"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Minute)
	defer cancel()


	err := os.MkdirAll(LogFilePath, os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stdout, constants.MsgCreateFolderFail)
	}

	logFullPath := filepath.Join(LogFilePath, LogFileName)

	file, err := os.OpenFile(logFullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stdout, constants.MsgCreateFileFail)
	}

	defer file.Close()

	slogHandler := slog.NewJSONHandler(file, nil)
	
	slog.SetDefault(slog.New(slogHandler))
	
	err = godotenv.Load()
	if err != nil {
		
		slog.Error(constants.MsgToolsInitFail, "tools", "godotenv", "error", err.Error())
		panic(err)
	}

	db, err := config.NewDB(ctx)
	if err != nil {
		slog.Error(constants.MsgToolsInitFail, "tools", "db", "error", err.Error())
		panic(err)
	}
	defer db.Close()

	repositories := repository.NewRepositories(db)
	services := service.NewServices(repositories)
	handlers := handler.NewHandlers(services)

	router := config.NewRouter(handlers)

	port := utils.GetEnv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	slog.Info("Server starting", "address", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		slog.Error("Server failed to start", "error", err.Error())
		panic(err)
	}
}