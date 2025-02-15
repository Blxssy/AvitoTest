package services_test

import (
	"context"
	"database/sql"
	"github.com/Blxssy/AvitoTest/internal/models"
	"github.com/Blxssy/AvitoTest/internal/repo"
	"github.com/Blxssy/AvitoTest/internal/repo/mocks"
	"github.com/Blxssy/AvitoTest/internal/services"
	mocks2 "github.com/Blxssy/AvitoTest/pkg/token/mocks"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestGetBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mocks.NewMockCoinRepository(ctrl)
	tokenGenMock := mocks2.NewMockTokenGenerator(ctrl)

	service := services.NewCoinService(repoMock, tokenGenMock)

	ctx := context.Background()
	tokenStr := "valid-token"
	username := "testuser"
	balance := 1000

	tokenGenMock.EXPECT().ParseToken(tokenStr).Return(username, nil)
	repoMock.EXPECT().GetBalance(ctx, gomock.Any()).Return(balance, nil)

	result, err := service.GetBalance(ctx, services.GetBalanceParams{Token: tokenStr})

	assert.NoError(t, err)
	assert.Equal(t, balance, result)
}

func TestAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mocks.NewMockCoinRepository(ctrl)
	tokenGenMock := mocks2.NewMockTokenGenerator(ctrl)

	service := services.NewCoinService(repoMock, tokenGenMock)

	ctx := context.Background()
	params := services.AuthParams{Username: "testuser", Password: "password"}

	repoMock.EXPECT().GetUserByUsername(ctx, params.Username).Return(nil, sql.ErrNoRows)
	repoMock.EXPECT().CreateUser(ctx, gomock.Any()).Return(nil)
	tokenGenMock.EXPECT().NewToken(params.Username).Return("new-token", nil)

	token, err := service.Auth(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, "new-token", token)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.MinCost)
	repoMock.EXPECT().GetUserByUsername(ctx, params.Username).Return(&models.User{Username: params.Username, PasswordHash: string(hashedPassword)}, nil)
	tokenGenMock.EXPECT().NewToken(params.Username).Return("auth-token", nil)

	token, err = service.Auth(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, "auth-token", token)
}

func TestSendCoins(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mocks.NewMockCoinRepository(ctrl)
	tokenGenMock := mocks2.NewMockTokenGenerator(ctrl)

	service := services.NewCoinService(repoMock, tokenGenMock)

	ctx := context.Background()
	params := services.TransactionParams{
		Token: "valid-token", ReceiverUsername: "receiver", Amount: 500,
	}
	senderUsername := "sender"

	tokenGenMock.EXPECT().ParseToken(params.Token).Return(senderUsername, nil)

	repoMock.EXPECT().GetUserByUsername(ctx, params.ReceiverUsername).
		Return(&models.User{Username: params.ReceiverUsername}, nil)

	repoMock.EXPECT().GetBalance(ctx, repo.GetBalanceParams{Username: senderUsername}).
		Return(1000, nil)

	tx := &sqlx.Tx{}
	repoMock.EXPECT().BeginTx(ctx).Return(tx, nil)
	repoMock.EXPECT().DecreaseBalance(ctx, tx, repo.ChangeBalanceParams{
		Username: senderUsername, Amount: params.Amount,
	}).Return(nil)
	repoMock.EXPECT().IncreaseBalance(ctx, tx, repo.ChangeBalanceParams{
		Username: params.ReceiverUsername, Amount: params.Amount,
	}).Return(nil)
	repoMock.EXPECT().SaveTransaction(ctx, repo.SaveTransactionParams{
		SenderUsername: senderUsername, ReceiverUsername: params.ReceiverUsername, Amount: params.Amount,
	}).Return(nil)

	repoMock.EXPECT().CommitTx(tx).Return(nil)

	err := service.SendCoins(ctx, params)
	assert.NoError(t, err)
}

func TestSendCoinsInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mocks.NewMockCoinRepository(ctrl)
	tokenGenMock := mocks2.NewMockTokenGenerator(ctrl)

	service := services.NewCoinService(repoMock, tokenGenMock)

	ctx := context.Background()
	params := services.GetTransactionsParams{Token: "valid-token"}
	username := "testuser"

	tokenGenMock.EXPECT().ParseToken(params.Token).Return(username, nil)
	repoMock.EXPECT().GetTransactions(ctx, username).Return([]models.Transaction{
		{SenderUsername: "testuser", ReceiverUsername: "receiver", Amount: 100},
	}, nil)

	transactions, err := service.SendCoinsInfo(ctx, params)
	assert.NoError(t, err)
	assert.Len(t, transactions, 1)
}

func TestReceivedCoinsInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mocks.NewMockCoinRepository(ctrl)
	tokenGenMock := mocks2.NewMockTokenGenerator(ctrl)

	service := services.NewCoinService(repoMock, tokenGenMock)

	ctx := context.Background()
	params := services.GetTransactionsParams{Token: "valid-token"}
	username := "testuser"

	tokenGenMock.EXPECT().ParseToken(params.Token).Return(username, nil)
	repoMock.EXPECT().ReceivedCoinsInfo(ctx, username).Return([]models.Transaction{
		{SenderUsername: "sender", ReceiverUsername: "testuser", Amount: 100},
	}, nil)

	transactions, err := service.ReceivedCoinsInfo(ctx, params)
	assert.NoError(t, err)
	assert.Len(t, transactions, 1)
}

func TestGetPurchases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mocks.NewMockCoinRepository(ctrl)
	tokenGenMock := mocks2.NewMockTokenGenerator(ctrl)

	service := services.NewCoinService(repoMock, tokenGenMock)

	ctx := context.Background()
	params := services.GetPurchasesParams{Token: "valid-token"}
	username := "testuser"

	tokenGenMock.EXPECT().ParseToken(params.Token).Return(username, nil)
	repoMock.EXPECT().GetPurchases(ctx, username).Return([]models.PurchaseItem{
		{Item: "item1", Count: 1},
	}, nil)

	purchases, err := service.GetPurchases(ctx, params)
	assert.NoError(t, err)
	assert.Len(t, purchases, 1)
}

func TestBuyItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mocks.NewMockCoinRepository(ctrl)
	tokenGenMock := mocks2.NewMockTokenGenerator(ctrl)

	service := services.NewCoinService(repoMock, tokenGenMock)

	ctx := context.Background()
	params := services.BuyItemParams{Token: "valid-token", Item: "item1"}
	username := "testuser"

	tokenGenMock.EXPECT().ParseToken(params.Token).Return(username, nil)
	repoMock.EXPECT().GetItem(ctx, params.Item).Return(models.Item{Name: "item1", Price: 100}, nil)
	repoMock.EXPECT().BuyItem(ctx, repo.BuyItemParams{Username: username, Item: "item1", Price: 100}).Return(nil)

	err := service.BuyItem(ctx, params)
	assert.NoError(t, err)
}
