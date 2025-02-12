package v1

import (
	"github.com/Blxssy/AvitoTest/internal/services"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handler struct {
	coinService services.CoinService
	logger      *zap.Logger
}

type HandlerConfig struct {
	CoinService services.CoinService
	Logger      *zap.Logger
}

func NewHandler(cfg HandlerConfig) *Handler {
	return &Handler{
		logger:      cfg.Logger,
		coinService: cfg.CoinService,
	}
}

func (h *Handler) Init(router fiber.Router) {
	h.initCoinRoutes(router)
}
