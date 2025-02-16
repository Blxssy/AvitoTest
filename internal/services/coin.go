package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Blxssy/AvitoTest/internal/models"
	"github.com/Blxssy/AvitoTest/internal/repo"
	"github.com/Blxssy/AvitoTest/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

var UnauthorizedError = errors.New("unauthorized")

type CoinService interface {
	GetBalance(ctx context.Context, params GetBalanceParams) (int, error)
	Auth(ctx context.Context, params AuthParams) (string, error)
	SendCoins(ctx context.Context, params TransactionParams) error
	SendCoinsInfo(ctx context.Context, params GetTransactionsParams) ([]models.Transaction, error)
	ReceivedCoinsInfo(ctx context.Context, params GetTransactionsParams) ([]models.Transaction, error)
	GetPurchases(ctx context.Context, params GetPurchasesParams) ([]models.PurchaseItem, error)
	BuyItem(ctx context.Context, params BuyItemParams) error
}

type coinService struct {
	repo     repo.CoinRepository
	tokenGen token.TokenGenerator
}

func NewCoinService(repo repo.CoinRepository, tg token.TokenGenerator) CoinService {
	return &coinService{
		repo:     repo,
		tokenGen: tg,
	}
}

func (s *coinService) GetBalance(ctx context.Context, params GetBalanceParams) (int, error) {
	username, err := s.tokenGen.ParseToken(params.Token)
	if err != nil {
		return 0, fmt.Errorf("s.tokenGen.ParseToken: %w", err)
	}

	balance, err := s.repo.GetBalance(ctx, repo.GetBalanceParams{
		Username: username,
	})
	if err != nil {
		return 0, fmt.Errorf("s.repo.GetBalance: %w", err)
	}

	return balance, nil
}

func (s *coinService) Auth(ctx context.Context, params AuthParams) (string, error) {
	user, err := s.repo.GetUserByUsername(ctx, params.Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("s.repo.GetUserByUsername: %w", err)
	}

	if user == nil {
		passHash, passErr := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.MinCost)
		if passErr != nil {
			return "", UnauthorizedError
		}

		err = s.repo.CreateUser(ctx, repo.CreateUserParams{
			Username: params.Username,
			PassHash: string(passHash),
			Balance:  1000,
		})
		if err != nil {
			return "", fmt.Errorf("s.repo.CreateUser: %w", err)
		}

		accessToken, tErr := s.tokenGen.NewToken(params.Username)
		if tErr != nil {
			return "", fmt.Errorf("s.tokenGen.NewToken: %w", err)
		}

		return accessToken, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password))
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}

	accessToken, err := s.tokenGen.NewToken(user.Username)
	if err != nil {
		return "", fmt.Errorf("s.tokenGen.NewToken: %w", err)
	}

	return accessToken, nil
}

func (s *coinService) SendCoins(ctx context.Context, params TransactionParams) error {
	senderUsername, err := s.tokenGen.ParseToken(params.Token)
	if err != nil {
		return fmt.Errorf("s.tokenGen.ParseToken: %w", err)
	}

	_, err = s.repo.GetUserByUsername(ctx, params.ReceiverUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("receiver not found")
		}
		return fmt.Errorf("s.repo.GetUserByUsername: %w", err)
	}

	balance, err := s.repo.GetBalance(ctx, repo.GetBalanceParams{
		Username: senderUsername,
	})
	if err != nil {
		return fmt.Errorf("s.repo.GetBalance: %w", err)
	}

	if balance < params.Amount {
		return fmt.Errorf("insufficient funds")
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("s.repo.BeginTx: %w", err)
	}

	defer func() {
		if err != nil {
			err = s.repo.RollbackTx(tx)
			if err != nil {
				err = fmt.Errorf("s.repo.RollbackTx: %w", err)
			}
		}
	}()

	if err = s.repo.DecreaseBalance(ctx, tx, repo.ChangeBalanceParams{
		Username: senderUsername, Amount: params.Amount,
	}); err != nil {
		return fmt.Errorf("s.repo.DecreaseBalance: %w", err)
	}

	if err = s.repo.IncreaseBalance(ctx, tx, repo.ChangeBalanceParams{
		Username: params.ReceiverUsername, Amount: params.Amount,
	}); err != nil {
		return fmt.Errorf("s.repo.IncreaseBalance: %w", err)
	}

	if err = s.repo.SaveTransaction(ctx, repo.SaveTransactionParams{
		SenderUsername: senderUsername, ReceiverUsername: params.ReceiverUsername, Amount: params.Amount,
	}); err != nil {
		return fmt.Errorf("s.repo.SaveTransaction: %w", err)
	}

	if err = s.repo.CommitTx(tx); err != nil {
		return fmt.Errorf("s.repo.CommitTx: %w", err)
	}

	return nil
}

func (s *coinService) SendCoinsInfo(ctx context.Context, params GetTransactionsParams) ([]models.Transaction, error) {
	username, err := s.tokenGen.ParseToken(params.Token)
	if err != nil {
		return nil, fmt.Errorf("s.tokenGen.ParseToken: %w", err)
	}

	transactions, err := s.repo.GetTransactions(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("s.repo.GetTransactions: %w", err)
	}

	return transactions, nil
}

func (s *coinService) ReceivedCoinsInfo(ctx context.Context, params GetTransactionsParams) ([]models.Transaction, error) {
	username, err := s.tokenGen.ParseToken(params.Token)
	if err != nil {
		return nil, fmt.Errorf("s.tokenGen.ParseToken: %w", err)
	}

	transactions, err := s.repo.ReceivedCoinsInfo(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("s.repo.ReceivedCoinsInfo: %w", err)
	}

	return transactions, nil
}

func (s *coinService) GetPurchases(ctx context.Context, params GetPurchasesParams) ([]models.PurchaseItem, error) {
	username, err := s.tokenGen.ParseToken(params.Token)
	if err != nil {
		return nil, fmt.Errorf("s.tokenGen.ParseToken: %w", err)
	}

	purchases, err := s.repo.GetPurchases(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("s.repo.GetPurchases: %w", err)
	}

	return purchases, nil
}

func (s *coinService) BuyItem(ctx context.Context, params BuyItemParams) error {
	username, err := s.tokenGen.ParseToken(params.Token)
	if err != nil {
		return fmt.Errorf("s.tokenGen.ParseToken: %w", err)
	}

	item, err := s.repo.GetItem(ctx, params.Item)
	if err != nil {
		return fmt.Errorf("s.repo.GetItem: %w", err)
	}

	err = s.repo.BuyItem(ctx, repo.BuyItemParams{
		Username: username,
		Item:     item.Name,
		Price:    item.Price,
	})
	if err != nil {
		return fmt.Errorf("s.repo.BuyItem: %w", err)
	}

	return nil
}
