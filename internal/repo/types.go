package repo

type GetBalanceParams struct {
	Username string
}

type GetUserByIDParams struct {
	UserID uint32
}

type GetUserByUsername struct {
	Username string
}

type CreateUserParams struct {
	Username string
	PassHash string
	Balance  int
}

type ChangeBalanceParams struct {
	Username string
	Amount   int
}

type SaveTransactionParams struct {
	SenderUsername   string
	ReceiverUsername string
	Amount           int
}

type GetTransactionsParams struct {
	Username string
}

type BuyItemParams struct {
	Username string
	Item     string
	Price    int
}
