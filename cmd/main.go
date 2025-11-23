package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"tg_bot/internal/adapters"
	"tg_bot/internal/config"
	"tg_bot/internal/service"

	"github.com/joho/godotenv"
)

func setupLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	return logger
}

func main() {
	logger := setupLogger()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := godotenv.Load(); err != nil {
		logger.Warn("Файл .env не найден или не загружен")
	}

	cfg, err := config.LoadConfig(logger)
	if err != nil {
		logger.Error("Ошибка загрузки конфигурации", "error", err)
		os.Exit(1)
	}

	quoteAPI := adapters.NewZenquotesAPI()
	mymemoryAPI := adapters.NewMymemoryAPI()
	weatherAPI := adapters.NewWeatherAPI(cfg.WeatherToken)
	valuteAPI := adapters.NewValuteAPI()

	quoteService := service.NewQuoteService(quoteAPI, mymemoryAPI, weatherAPI, valuteAPI)
	telegramAPI, err := adapters.NewTelegramAdapter(cfg.BotToken)
	if err != nil {
		logger.Error("ошибка создания адаптера telegramAPI")
	}

	telegramService := service.NewSendMessageService(telegramAPI, quoteService)
	telegramService.StartListening(ctx, telegramAPI)

	logger.Info("Бот запущен и слушает сообщения...")
	<-ctx.Done()
	logger.Info("Получен сигнал завершения. Останавливаем приложение...")

}
