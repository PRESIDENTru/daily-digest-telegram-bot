package service

import (
	"context"
	"errors"
	"testing"
	"tg_bot/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockQuotesAPI struct {
	mock.Mock
}

func (m *MockQuotesAPI) GetQuote(ctx context.Context) (*models.Quote, error) {
	args := m.Called(ctx)
	//?????
	return args.Get(0).(*models.Quote), args.Error(1)
}

type MockWeatherAPI struct {
	mock.Mock
}

func (m *MockWeatherAPI) GetWeather(ctx context.Context, city string) (*models.WeatherResponse, error) {
	args := m.Called(ctx, city)
	return args.Get(0).(*models.WeatherResponse), args.Error(1)
}

type MockValuteAPI struct {
	mock.Mock
}

func (m *MockValuteAPI) GetValute(ctx context.Context, base string) (*models.Valute, error) {
	args := m.Called(ctx, base)
	return args.Get(0).(*models.Valute), args.Error(1)
}

type MockTranlateAPI struct {
	mock.Mock
}

func (m *MockTranlateAPI) GetTranslate(ctx context.Context, text string) (*models.Translation, error) {
	args := m.Called(ctx, text)
	return args.Get(0).(*models.Translation), args.Error(1)
}

func TestNewDataAggreratorService(t *testing.T) {
	mockQuoteAPI := new(MockQuotesAPI)
	mockTranslateAPI := new(MockTranlateAPI)
	mockWeatherAPI := new(MockWeatherAPI)
	MockValuteAPI := new(MockValuteAPI)

	service := NewDataAggregatorService(
		mockQuoteAPI,
		mockTranslateAPI,
		mockWeatherAPI,
		MockValuteAPI,
	)

	assert.NotNil(t, service)
}

// Тест успешной работы GetWeather
func TestDataAggregator_GetWeatherInfo_Success(t *testing.T) {
	mockWeatherAPI := new(MockWeatherAPI)
	aggregator := NewDataAggregatorService(nil, nil, mockWeatherAPI, nil)

	expectedWeather := &models.WeatherResponse{
		Location: struct {
			Name    string `json:"name"`
			Region  string `json:"region"`
			Country string `json:"country"`
		}{
			Name:    "Moscow",
			Country: "Russia",
		},
		Current: struct {
			TempC float64 `json:"temp_c"`
			TempF float64 `json:"temp_f"`
		}{
			TempC: 20.5,
		},
	}
	mockWeatherAPI.On("GetWeather", mock.Anything, "Moscow").Return(expectedWeather, nil)

	result, err := aggregator.GetWeatherInfo(context.Background(), "Moscow")

	assert.NoError(t, err)
	assert.Contains(t, result, "Moscow")
	assert.Contains(t, result, "20.5")
	mockWeatherAPI.AssertCalled(t, "GetWeather", mock.Anything, "Moscow")
}

func TestDataAggregator_GetWeatherInfo_Fail(t *testing.T) {
	mockWeatherAPI := new(MockWeatherAPI)
	aggregator := NewDataAggregatorService(nil, nil, mockWeatherAPI, nil)

	expectedWeather := &models.WeatherResponse{
		Location: struct {
			Name    string `json:"name"`
			Region  string `json:"region"`
			Country string `json:"country"`
		}{
			Name:    "Moscow",
			Country: "Russia",
		},
		Current: struct {
			TempC float64 `json:"temp_c"`
			TempF float64 `json:"temp_f"`
		}{
			TempC: 20.5,
		},
	}
	mockWeatherAPI.On("GetWeather", mock.Anything, "Moscow").Return(expectedWeather, errors.New("не дам данные!"))

	_, err := aggregator.GetWeatherInfo(context.Background(), "Moscow")

	assert.NoError(t, err)
	mockWeatherAPI.AssertCalled(t, "GetWeather", mock.Anything, "Moscow")
}
