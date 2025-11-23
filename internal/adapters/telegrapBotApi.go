package adapters

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramAdapter struct {
	bot     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
}

func NewTelegramAdapter(botToken string) (*TelegramAdapter, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	return &TelegramAdapter{
		bot:     bot,
		updates: updates,
	}, nil
}

func (t *TelegramAdapter) SendMessage(ctx context.Context, chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)

	t.bot.Debug = true
	_, err := t.bot.Send(msg)
	return err
}

func (t *TelegramAdapter) StartListening(ctx context.Context, handler func(update tgbotapi.Update)) {
	go func() {
		for update := range t.updates {
			handler(update)
		}
	}()
}
