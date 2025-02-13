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
		coinRoute.Get("info", h.Info)
		coinRoute.Get("buy/:item", h.BuyItem)
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

type SendCoinRequest struct {
	ReceiverUsername string `json:"toUser"`
	Amount           int    `json:"amount"`
}

func (h *Handler) Transaction(ctx *fiber.Ctx) error {
	var req SendCoinRequest
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

	err = h.coinService.SendCoins(ctx.Context(), services.TransactionParams{
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

func (h *Handler) Info(ctx *fiber.Ctx) error {
	token, err := getToken(ctx)
	if err != nil {
		return err
	}

	balance, err := h.coinService.GetBalance(ctx.Context(), services.GetBalanceParams{
		Token: token,
	})
	if err != nil {
		if errors.Is(err, services.UnauthorizedError) {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.coinService.GetBalance: %v", err))
	}

	purchases, err := h.coinService.GetPurchases(ctx.Context(), services.GetPurchasesParams{
		Token: token,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.coinService.GetPurchases: %v", err))
	}

	fInventory := make([]fiber.Map, len(purchases))
	for i, p := range purchases {
		fInventory[i] = fiber.Map{
			"type":     p.Item,
			"quantity": p.Count,
		}
	}

	sentCoins, err := h.coinService.SendCoinsInfo(ctx.Context(), services.GetTransactionsParams{
		Token: token,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.coinService.GetTransactions: %v", err))
	}

	fSentCoins := make([]fiber.Map, len(sentCoins))
	for i, transaction := range sentCoins {
		fSentCoins[i] = fiber.Map{
			"toUser": transaction.ReceiverUsername,
			"amount": transaction.Amount,
		}
	}

	receivedCoins, err := h.coinService.ReceivedCoinsInfo(ctx.Context(), services.GetTransactionsParams{
		Token: token,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.coinService.ReceivedCoinsInfo: %v", err))
	}

	fReceivedCoins := make([]fiber.Map, len(receivedCoins))
	for i, t := range receivedCoins {
		fReceivedCoins[i] = fiber.Map{
			"fromUser": t.SenderUsername,
			"amount":   t.Amount,
		}
	}

	return ctx.JSON(fiber.Map{
		"coins":     balance,
		"inventory": fInventory,
		"coinHistory": fiber.Map{
			"received": fReceivedCoins,
			"sent":     fSentCoins,
		},
	})
}

func (h *Handler) BuyItem(ctx *fiber.Ctx) error {
	token, err := getToken(ctx)
	if err != nil {
		return err
	}
	item := ctx.Params("item")

	err = h.coinService.BuyItem(ctx.Context(), services.BuyItemParams{
		Token: token,
		Item:  item,
	})
	if err != nil {
		if errors.Is(err, services.UnauthorizedError) {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("h.coinService.BuyItem: %v", err))
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
