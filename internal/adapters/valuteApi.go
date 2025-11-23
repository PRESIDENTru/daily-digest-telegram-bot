package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"tg_bot/internal/models"
	"time"
)

type ValuteAPI struct {
	client *http.Client
}

func NewValuteAPI() *ValuteAPI {
	return &ValuteAPI{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (v *ValuteAPI) GetValute(ctx context.Context, baseCurrency string) (*models.Valute, error) {
	url := fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/%s", baseCurrency)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, errors.New("exchangeRateAPI: ошибка создания запроса")
	}

	resp, err := v.client.Do(req)
	if err != nil {
		return nil, errors.New("exchangeRateAPI: ошибка запроса к API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("exchangeRateAPI: статус ошибки %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("exchangeRateAPI: ошибка чтения body")
	}

	var valute models.Valute

	if err := json.Unmarshal(body, &valute); err != nil {
		return nil, errors.New("exchangeRateAPI: ошибка unmarshal ответа")
	}

	return &valute, nil
}
