package pg

import (
	"context"
	"errors"
	"github.com/Blxssy/AvitoTest/internal/models"
	"github.com/Blxssy/AvitoTest/internal/repo"
)

type Purchase struct {
	Username string `db:"username"`
	Item     string `db:"item"`
	Price    int    `db:"price"`
}

type Item struct {
	Name  string
	Price int
}

const repoStmtGetPurchases = `
SELECT item, COUNT(*) as count
FROM purchases
WHERE username = $1
GROUP BY item
`

const repoStmtBuyItem = `
INSERT INTO purchases (username, item, price)
VALUES ($1, $2, $3)
`

const repoStmtGetItems = `
SELECT name, price
FROM items
WHERE name = $1
`

func (r *CoinRepo) GetPurchases(ctx context.Context, username string) ([]models.PurchaseItem, error) {
	rows, err := r.db.QueryxContext(ctx, repoStmtGetPurchases, username)
	if err != nil {
		return nil, err
	}

	var purchases []models.PurchaseItem
	for rows.Next() {
		var purchase models.PurchaseItem
		if err := rows.StructScan(&purchase); err != nil {
			return nil, err
		}

		purchases = append(purchases, purchase)
	}

	return purchases, nil
}

func (r *CoinRepo) BuyItem(ctx context.Context, params repo.BuyItemParams) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var balance int
	err = tx.QueryRowContext(ctx, "SELECT balance FROM users WHERE username=$1", params.Username).Scan(&balance)
	if err != nil {
		return err
	}

	if balance < params.Price {
		return errors.New("insufficient funds")
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET balance = balance - $1 WHERE username = $2", params.Price, params.Username)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, repoStmtBuyItem,
		params.Username, params.Item, params.Price)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *CoinRepo) GetItem(ctx context.Context, itemName string) (models.Item, error) {
	var item models.Item
	err := r.db.QueryRowxContext(ctx, repoStmtGetItems, itemName).StructScan(&item)
	if err != nil {
		return models.Item{}, err
	}

	return item, nil
}
