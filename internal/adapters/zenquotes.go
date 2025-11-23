package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"tg_bot/internal/models"
	"time"
)

type ZenquotesAPI struct {
	client *http.Client
}

func NewZenquotesAPI() *ZenquotesAPI {
	return &ZenquotesAPI{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetQuote получает случайную фразу на английском языке из Zenquotes API
// Возвращает структуру Quote или ошибку, если запрос или декодирование не удалось выполнить
func (z *ZenquotesAPI) GetQuote(ctx context.Context) (*models.Quote, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://zenquotes.io/api/random", nil)

	if err != nil {
		return nil, errors.New("zqnquotes: ошибка создания запроса")
	}

	resp, err := z.client.Do(req)
	if err != nil {
		return nil, errors.New("ошибка запроса к API")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("zqnquotes: ошибка чтения body")
	}
	var quote []models.Quote

	if err := json.Unmarshal(body, &quote); err != nil {
		log.Fatal(err)
		return nil, errors.New("zqnquotes: ошибка unmarshal")
	}

	if quote[0].Text == "" {
		return nil, errors.New("получена пустая цитата")
	}

	author := quote[0].Autor
	if author == "" {
		author = "Неизвестный автор"
	}

	return &quote[0], nil
}
