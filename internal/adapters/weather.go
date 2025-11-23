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

type WeatherAPI struct {
	client *http.Client
	token  string
}

func NewWeatherAPI(token string) *WeatherAPI {
	return &WeatherAPI{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		token: token,
	}
}

func (w *WeatherAPI) GetWeather(ctx context.Context, city string) (*models.WeatherResponse, error) {
	//валидацию города
	req_str := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&lang=ru", w.token, city)
	req, err := http.NewRequestWithContext(ctx, "GET", req_str, nil)

	if err != nil {
		return nil, errors.New("weatherAPI: ошибка создания запроса")
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.New("weatherAPI: ошибка запроса к API")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("weatherAPI: ошибка чтения body")
	}
	var weather models.WeatherResponse
	if err := json.Unmarshal(body, &weather); err != nil {
		return nil, errors.New("weatherAPI: ошибка unmarshal ответа")
	}

	return &weather, nil
}
