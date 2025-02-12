package pg

import (
	"context"
	"fmt"
	"github.com/Blxssy/AvitoTest/internal/models"
	"github.com/Blxssy/AvitoTest/internal/repo"
	"github.com/jmoiron/sqlx"
)

type CoinRepo struct {
	db *sqlx.DB
}

func NewCoinRepo(db *sqlx.DB) *CoinRepo {
	return &CoinRepo{db: db}
}

type User struct {
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
	Balance      int    `db:"balance"`
}

const repoStmtFindByUsername = `
select * 
from users 
where username = $1
`

const repoStmtGetBalance = `
select * 
from users
where id = $1
`

const repoStmtCreateUser = `
insert into 
    users
    (username, password_hash, balance)
    values ($1, $2, $3);
`

const repoStmtDecreaseBalance = `
UPDATE users 
SET balance = balance - $1 
WHERE username = $2
`

const repoStmtIncreaseBalance = `
UPDATE users 
SET balance = balance + $1 
WHERE username = $2
`

func (r *CoinRepo) GetBalance(ctx context.Context, params repo.GetBalanceParams) (int, error) {
	var balance int
	if err := r.db.SelectContext(
		ctx,
		&balance,
		repoStmtGetBalance,
		params.UserID,
	); err != nil {
		return 0, fmt.Errorf("r.db.SelectContext: %w", err)
	}
	return balance, nil
}

func (r *CoinRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var usr User
	if err := r.db.GetContext(
		ctx,
		&usr,
		repoStmtFindByUsername,
		username,
	); err != nil {
		return nil, fmt.Errorf("r.db.GetContext: %w", err)
	}
	return &models.User{
		Username:     usr.Username,
		PasswordHash: usr.PasswordHash,
		Balance:      usr.Balance,
	}, nil
}

func (r *CoinRepo) CreateUser(ctx context.Context, params repo.CreateUserParams) error {
	if _, err := r.db.ExecContext(
		ctx,
		repoStmtCreateUser,
		params.Username,
		params.PassHash,
		params.Balance,
	); err != nil {
		return fmt.Errorf("r.db.ExecContext: %w", err)
	}
	return nil
}

func (r *CoinRepo) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}

func (r *CoinRepo) DecreaseBalance(ctx context.Context, tx *sqlx.Tx, params repo.ChangeBalanceParams) error {
	_, err := tx.ExecContext(
		ctx,
		repoStmtDecreaseBalance,
		params.Amount,
		params.Username,
	)
	if err != nil {
		return fmt.Errorf("tx.ExecContext (decrease): %w", err)
	}
	return nil
}

func (r *CoinRepo) IncreaseBalance(ctx context.Context, tx *sqlx.Tx, params repo.ChangeBalanceParams) error {
	_, err := tx.ExecContext(
		ctx,
		repoStmtIncreaseBalance,
		params.Amount,
		params.Username,
	)
	if err != nil {
		return fmt.Errorf("tx.ExecContext (increase): %w", err)
	}
	return nil
}
