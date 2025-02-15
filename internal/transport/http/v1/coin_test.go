package v1_test

import (
	"github.com/Blxssy/AvitoTest/internal/services"
	"github.com/Blxssy/AvitoTest/internal/services/mocks"
	"github.com/Blxssy/AvitoTest/internal/transport/http/v1"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockCoinService(ctrl)

	app := fiber.New()
	handler := v1.NewHandler(v1.HandlerConfig{
		CoinService: mockService,
		Logger:      nil,
	})
	handler.Init(app)
	app.Post("/api/auth", handler.Auth)

	requestBody := `{"username":"test","password":"123"}`
	mockService.EXPECT().Auth(gomock.Any(), services.AuthParams{
		Username: "test",
		Password: "123",
	}).Return("valid_token", nil)

	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/auth", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestBuyItemHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockCoinService(ctrl)

	app := fiber.New()
	handler := v1.NewHandler(v1.HandlerConfig{
		CoinService: mockService,
		Logger:      nil,
	})
	handler.Init(app)
	app.Get("/api/buy/:item", handler.BuyItem)

	mockService.EXPECT().BuyItem(gomock.Any(), services.BuyItemParams{
		Token: "valid_token",
		Item:  "powerbank",
	}).Return(nil)

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/buy/powerbank", nil) // Используем реальное значение
	req.Header.Set("Authorization", "Bearer valid_token")

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSendCoinsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockCoinService(ctrl)

	app := fiber.New()
	handler := v1.NewHandler(v1.HandlerConfig{
		CoinService: mockService,
		Logger:      nil,
	})
	handler.Init(app)
	app.Post("/api/sendCoin", handler.Transaction)

	requestBody := `{"toUser": "Bill", "amount": 100}`
	mockService.EXPECT().SendCoins(gomock.Any(), services.TransactionParams{
		Token:            "valid_token",
		ReceiverUsername: "Bill",
		Amount:           100,
	}).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/sendCoin", strings.NewReader(requestBody))
	req.Header.Set("Authorization", "Bearer valid_token")
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
