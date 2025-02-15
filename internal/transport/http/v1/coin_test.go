package v1_test

import (
	"bytes"
	"encoding/json"
	"github.com/Blxssy/AvitoTest/internal/transport/http/v1"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupApp() *fiber.App {
	app := fiber.New()
	handler := v1.NewHandler(v1.HandlerConfig{})
	handler.Init(app)
	return app
}

func TestBuyItem(t *testing.T) {
	app := setupApp()

	// Получение токена
	token := "your_test_token"

	// Покупка мерча
	req := httptest.NewRequest(http.MethodGet, "/api/buy/item1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSendCoin(t *testing.T) {
	app := setupApp()

	// Получение токена
	token := "your_test_token"

	// Передача монеток
	sendCoinReq := v1.SendCoinRequest{
		ReceiverUsername: "receiveruser",
		Amount:           10,
	}
	sendCoinReqBody, _ := json.Marshal(sendCoinReq)
	req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(sendCoinReqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
