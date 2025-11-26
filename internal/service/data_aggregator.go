package service

import (
	"context"
	"fmt"
	"tg_bot/internal/interfaces"
	"tg_bot/internal/models"
)

type DataAggregatorService struct {
	quoteAPI     interfaces.QouteAPI
	translateAPI interfaces.TranlateAPI
	weatherAPI   interfaces.WeatherAPI
	valuteAPI    interfaces.ValuteAPI
}

func NewDataAggregatorService(
	quoteAPI interfaces.QouteAPI,
	translateAPI interfaces.TranlateAPI,
	weatherAPI interfaces.WeatherAPI,
	valuteAPI interfaces.ValuteAPI,
) *DataAggregatorService {
	return &DataAggregatorService{
		quoteAPI:     quoteAPI,
		translateAPI: translateAPI,
		weatherAPI:   weatherAPI,
		valuteAPI:    valuteAPI,
	}
}

func (s *DataAggregatorService) GetTranslatedQuote(ctx context.Context) (string, error) {
	quote, err := s.quoteAPI.GetQuote(ctx)

	if err != nil {
		return "", fmt.Errorf("ошибка получения цитаты: %w", err)
	}

	translated, err := s.translateAPI.GetTranslate(ctx, quote.Text)
	if err != nil {
		return "", fmt.Errorf("ошибка перевода цитаты: %w", err)
	}
	result := fmt.Sprintf("💬 %s \nАвтор: %s", translated.TextData.Text, quote.Autor)
	return result, nil
}

func (s *DataAggregatorService) GetWeatherInfo(ctx context.Context, city string) (string, error) {
	weather, err := s.weatherAPI.GetWeather(ctx, city)

	if err != nil {
		return "", fmt.Errorf("ошибка получения информации о погоде: %w", err)
	}

	result := fmt.Sprintf("☀️ Страна: %s, Город: %s, Температура: %.1f°C\n",
		weather.Location.Country,
		weather.Location.Name,
		weather.Current.TempC,
	)
	return result, nil
}

func (s *DataAggregatorService) GetValuteRUB(ctx context.Context) (string, error) {
	results := make(chan *models.Valute, 3)
	go func() {
		usd, err := s.valuteAPI.GetValute(ctx, "USD")
		if err != nil {
			results <- &models.Valute{
				Rates: make(map[string]float64),
				Err:   err,
				Code:  "USD",
			}
		} else {
			results <- &models.Valute{
				Rates: usd.Rates,
				Err:   err,
				Code:  "USD",
			}
		}
	}()
	go func() {
		cny, err := s.valuteAPI.GetValute(ctx, "CNY")
		if err != nil {
			results <- &models.Valute{
				Rates: make(map[string]float64),
				Err:   err,
				Code:  "CNY",
			}
		} else {
			results <- &models.Valute{
				Rates: cny.Rates,
				Err:   err,
				Code:  "CNY",
			}
		}

	}()
	go func() {
		eur, err := s.valuteAPI.GetValute(ctx, "EUR")
		if err != nil {
			results <- &models.Valute{
				Rates: make(map[string]float64),
				Err:   err,
				Code:  "EUR",
			}
		} else {
			results <- &models.Valute{
				Rates: eur.Rates,
				Err:   err,
				Code:  "EUR",
			}
		}
	}()

	valutes := make(map[string]map[string]float64)
	var errors []error

	for i := 0; i < 3; i++ {
		result := <-results
		if result.Err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", result.Code, result.Err))
		} else {
			valutes[result.Code] = result.Rates
		}
	}

	if len(errors) > 0 {
		return "", fmt.Errorf("ошибки получения курсов: %v", errors)
	}

	result := fmt.Sprintf("💵1 доллар = %0.2f руб\n💴1 юань = %0.2fруб\n💶1 евро = %0.2f руб\n",
		valutes["USD"]["RUB"],
		valutes["CNY"]["RUB"],
		valutes["EUR"]["RUB"],
	)
	return result, nil
}
