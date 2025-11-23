package config

import (
	"errors"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken     string
	WeatherToken string
}

func LoadConfig(logger *slog.Logger) (*Config, error) {
	botToken := os.Getenv("BOT_TOKEN")
	weatherToken := os.Getenv("WEATHER_TOKEN")

	if botToken == "" {
		logger.Error("Переменные окружения отсутсвуют")
		if err := godotenv.Load(); err != nil {
			logger.Error("Не удалось загрузить .env файл")
			return nil, errors.New("файл .env не найден")
		}
	}

	if botToken == "" {
		logger.Error("BOT_TOKEN is required")
		return nil, errors.New("BOT_TOKEN environment variable is required")
	}

	return &Config{
		BotToken:     botToken,
		WeatherToken: weatherToken,
	}, nil
}
