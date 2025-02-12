package v1

import (
	"errors"
	"fmt"
	"github.com/Blxssy/AvitoTest/internal/services"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func (h *Handler) initCoinRoutes(router fiber.Router) {
	coinRoute := router.Group("/api")
	_ = coinRoute
	{
		coinRoute.Post("auth", h.Auth)
		coinRoute.Post("sendCoin", h.Transaction)
	}
}

type AuthRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (h *Handler) Auth(ctx *fiber.Ctx) error {
	var req AuthRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			fmt.Errorf("ctx.BodyParser: %w", err).Error(),
		)
	}

	accessToken, err := h.coinService.Auth(ctx.Context(), services.AuthParams{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, services.UnauthorizedError) {
			return fiber.NewError(
				fiber.StatusUnauthorized,
				fmt.Errorf("unauthorized").Error())
		}
		return fiber.NewError(
			fiber.StatusInternalServerError,
			fmt.Errorf("h.coinService.Auth: %w", err).Error(),
		)
	}

	return ctx.JSON(fiber.Map{
		"token": accessToken,
	})
}

type TransactionRequest struct {
	ReceiverUsername string `json:"receiver_username"`
	Amount           int    `json:"amount"`
}

func (h *Handler) Transaction(ctx *fiber.Ctx) error {
	var req TransactionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			fmt.Errorf("ctx.BodyParser: %w", err).Error(),
		)
	}

	token, err := getToken(ctx)
	if err != nil {
		return err
	}

	err = h.coinService.Transaction(ctx.Context(), services.TransactionParams{
		Token:            token,
		ReceiverUsername: req.ReceiverUsername,
		Amount:           req.Amount,
	})
	if err != nil {
		if errors.Is(err, services.UnauthorizedError) {
			return fiber.NewError(
				fiber.StatusUnauthorized,
				fmt.Errorf("unauthorized").Error(),
			)
		}
		return fiber.NewError(
			fiber.StatusInternalServerError,
			fmt.Errorf("h.coinService.Transaction: %w", err).Error(),
		)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func getToken(ctx *fiber.Ctx) (string, error) {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		return "", fiber.NewError(fiber.StatusUnauthorized, "invalid authorization header")
	}

	return token, nil
}
