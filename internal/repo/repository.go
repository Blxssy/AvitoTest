package repo

import (
	"context"
	"github.com/Blxssy/AvitoTest/internal/models"
	"github.com/jmoiron/sqlx"
)

type CoinRepository interface {
	GetBalance(ctx context.Context, params GetBalanceParams) (int, error)
	CreateUser(ctx context.Context, params CreateUserParams) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
	DecreaseBalance(ctx context.Context, tx *sqlx.Tx, params ChangeBalanceParams) error
	IncreaseBalance(ctx context.Context, tx *sqlx.Tx, params ChangeBalanceParams) error
	SaveTransaction(ctx context.Context, params SaveTransactionParams) error
	GetTransactions(ctx context.Context, username string) ([]models.Transaction, error)
	ReceivedCoinsInfo(ctx context.Context, username string) ([]models.Transaction, error)
	GetPurchases(ctx context.Context, username string) ([]models.PurchaseItem, error)
	BuyItem(ctx context.Context, params BuyItemParams) error
	GetItem(ctx context.Context, itemName string) (models.Item, error)
	CommitTx(tx *sqlx.Tx) error
	RollbackTx(tx *sqlx.Tx) error
}
