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
	//City Для погоды
}

func LoadConfig(logger *slog.Logger) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logger.Error("Не удалось загрузить .env файл")
		return nil, errors.New("файл .env не найден")
	}
	botToken := os.Getenv("BOT_TOKEN")
	weatherToken := os.Getenv("WEATHER_TOKEN")

	if botToken == "" {
		logger.Error("Переменные окружения отсутсвуют")
		return nil, errors.New("необходимые переменные окружения отсутсвуют")
	}

	//Валидация

	return &Config{
		BotToken:     botToken,
		WeatherToken: weatherToken,
	}, nil
}
