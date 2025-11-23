package service

import (
	"context"
	"log"
	"tg_bot/internal/interfaces"
	"tg_bot/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SendMessageService struct {
	telegram interfaces.TelegramSender
	quote    interfaces.QuoteService
}

func NewSendMessageService(telegram interfaces.TelegramSender, quote interfaces.QuoteService) *SendMessageService {
	return &SendMessageService{
		telegram: telegram,
		quote:    quote,
	}
}

func (s *SendMessageService) StartListening(ctx context.Context, receiver interfaces.TelegramReceiver) {
	receiver.StartListening(ctx, func(update tgbotapi.Update) {
		s.HandleUpdate(update, ctx)
	})
}

// city из чата либо базово из конфига
func (s *SendMessageService) HandleUpdate(update tgbotapi.Update, ctx context.Context) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	text := update.Message.Text

	log.Printf("Сообщение от %d: %s", userID, text)

	results := make(chan *models.Response, 3)
	go func() {
		quote, err := s.quote.GetTranslatedQuote(ctx)
		results <- &models.Response{
			Result: quote,
			Err:    err,
		}
	}()
	go func() {
		weather, err := s.quote.GetWeatherInfo(ctx, "Vladivostok")
		results <- &models.Response{
			Result: weather,
			Err:    err,
		}
	}()
	go func() {
		valute, err := s.quote.GetValuteRUB(ctx)
		results <- &models.Response{
			Result: valute,
			Err:    err,
		}
	}()

	var values []string
	var errors []error

	for i := 0; i < 3; i++ {
		result := <-results
		if result.Err != nil {
			errors = append(errors, result.Err)
			values = append(values, "Информация временно недоступна")
		} else {
			values = append(values, result.Result)
		}
	}
	if len(errors) > 0 {
		log.Println(errors)
	}

	message := values[0] + "\n" + values[1] + "\n" + values[2]
	if err := s.telegram.SendMessage(ctx, chatID, message); err != nil {
		log.Printf("Ошибка отправки цитаты: %v", err)
	}

}
