package interfaces

import (
	"context"
	"tg_bot/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type QouteAPI interface {
	GetQuote(ctx context.Context) (*models.Quote, error)
}

type TranlateAPI interface {
	GetTranslate(ctx context.Context, message string) (*models.Translation, error)
}

type WeatherAPI interface {
	GetWeather(ctx context.Context, city string) (*models.WeatherResponse, error)
}

type ValuteAPI interface {
	GetValute(ctx context.Context, baseCurrency string) (*models.Valute, error)
}

type TelegramSender interface {
	SendMessage(ctx context.Context, chatID int64, message string) error
}

type TelegramReceiver interface {
	StartListening(ctx context.Context, handler func(update tgbotapi.Update))
}

type QuoteService interface {
	GetTranslatedQuote(ctx context.Context) (string, error)
	GetWeatherInfo(ctx context.Context, city string) (string, error)
	GetValuteRUB(ctx context.Context) (string, error)
}
