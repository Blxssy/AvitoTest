package pg

import (
	"context"
	"github.com/Blxssy/AvitoTest/internal/models"
	"github.com/Blxssy/AvitoTest/internal/repo"
	"time"
)

type Transaction struct {
	ID               uint32    `db:"id"`
	SenderUsername   string    `db:"sender_username"`
	ReceiverUsername string    `db:"receiver_username"`
	Amount           int       `db:"amount"`
	CreatedAt        time.Time `db:"created_at"`
}

const repoStmtSaveTransaction = `
insert into
transactions
(sender_username, receiver_username, amount)
values ($1, $2, $3);
`

const repoStmtGetTransactions = `
select *
from transactions
where sender_username = $1
`

const repoStmtReceivedCoins = `
select *
from transactions
where receiver_username = $1
`

func (r *CoinRepo) SaveTransaction(ctx context.Context, params repo.SaveTransactionParams) error {
	_, err := r.db.ExecContext(
		ctx,
		repoStmtSaveTransaction,
		params.SenderUsername,
		params.ReceiverUsername,
		params.Amount,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *CoinRepo) GetTransactions(ctx context.Context, username string) ([]models.Transaction, error) {
	rows, err := r.db.QueryxContext(
		ctx,
		repoStmtGetTransactions,
		username,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction Transaction
		if err := rows.StructScan(&transaction); err != nil {
			return nil, err
		}

		transactions = append(transactions, models.Transaction{
			ID:               transaction.ID,
			SenderUsername:   transaction.SenderUsername,
			ReceiverUsername: transaction.ReceiverUsername,
			Amount:           transaction.Amount,
			CreatedAt:        transaction.CreatedAt,
		})
	}

	return transactions, nil
}

func (r *CoinRepo) ReceivedCoinsInfo(ctx context.Context, username string) ([]models.Transaction, error) {
	rows, err := r.db.QueryxContext(
		ctx,
		repoStmtReceivedCoins,
		username,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction Transaction
		if err := rows.StructScan(&transaction); err != nil {
			return nil, err
		}

		transactions = append(transactions, models.Transaction{
			ID:               transaction.ID,
			SenderUsername:   transaction.SenderUsername,
			ReceiverUsername: transaction.ReceiverUsername,
			Amount:           transaction.Amount,
			CreatedAt:        transaction.CreatedAt,
		})
	}

	return transactions, nil
}
