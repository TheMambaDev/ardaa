package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

var LogFile *os.File

func Setup() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("CONFIG", "error: ", err)
	}

	LogFile, err := os.OpenFile("./logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		slog.Error("CONFIG", "error: ", err)
	}

	// Create a custom json logger
	logger := slog.New(slog.NewJSONHandler(LogFile, nil))
	slog.SetDefault(logger)
}
