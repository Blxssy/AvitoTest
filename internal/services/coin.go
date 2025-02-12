package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Blxssy/AvitoTest/internal/repo"
	"github.com/Blxssy/AvitoTest/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

var UnauthorizedError = errors.New("unauthorized")

type CoinService interface {
	GetBalance(ctx context.Context, params GetBalanceParams) (int, error)
	Auth(ctx context.Context, params AuthParams) (string, error)
	Transaction(ctx context.Context, params TransactionParams) error
}

type coinService struct {
	repo     repo.CoinRepository
	tokenGen *token.TokenGen
}

func NewCoinService(repo repo.CoinRepository, tg *token.TokenGen) CoinService {
	return &coinService{
		repo:     repo,
		tokenGen: tg,
	}
}

func (s *coinService) GetBalance(ctx context.Context, params GetBalanceParams) (int, error) {
	balance, err := s.repo.GetBalance(ctx, repo.GetBalanceParams{
		UserID: params.UserID,
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
		passHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
		if err != nil {
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

		accessToken, err := s.tokenGen.NewToken(params.Username)
		if err != nil {
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

func (s *coinService) Transaction(ctx context.Context, params TransactionParams) error {
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

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("s.repo.BeginTx: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			return
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = s.repo.DecreaseBalance(ctx, tx, repo.ChangeBalanceParams{
		Username: senderUsername,
		Amount:   params.Amount,
	})
	if err != nil {
		return fmt.Errorf("s.repo.DecreaseBalance: %w", err)
	}

	err = s.repo.IncreaseBalance(ctx, tx, repo.ChangeBalanceParams{
		Username: params.ReceiverUsername,
		Amount:   params.Amount,
	})
	if err != nil {
		return fmt.Errorf("s.repo.IncreaseBalance: %w", err)
	}

	return nil
}
