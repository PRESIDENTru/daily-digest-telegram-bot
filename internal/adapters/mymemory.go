package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"tg_bot/internal/models"
	"time"
)

type MymemoryAPI struct {
	client *http.Client
}

func NewMymemoryAPI() *MymemoryAPI {
	return &MymemoryAPI{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetTranslate переводи переданную фразу originalText на русский язык
// Возвращает структуру Translation или ошибку, если запрос или декодирование не удалось выполнить
func (m *MymemoryAPI) GetTranslate(ctx context.Context, originalText string) (*models.Translation, error) {
	originalText = strings.ReplaceAll(originalText, " ", "%20")
	request_str := fmt.Sprintf("https://api.mymemory.translated.net/get?q=%s&langpair=en|ru", originalText)
	req, err := http.NewRequestWithContext(ctx, "GET", request_str, nil)
	if err != nil {
		return nil, errors.New("mymemory: ошибка создания запроса")
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, errors.New("mymemory: ошибка создания запроса")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("mymemory: ошибка чтения body")
	}

	var tranlate models.Translation

	if err := json.Unmarshal(body, &tranlate); err != nil {
		log.Fatal(err)
		return nil, errors.New("mymemory: ошибка unmarshal")
	}

	return &tranlate, nil
}
