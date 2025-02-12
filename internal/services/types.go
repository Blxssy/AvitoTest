package services

type GetBalanceParams struct {
	UserID uint32
}

type AuthParams struct {
	Username string
	Password string
}

type TransactionParams struct {
	Token            string
	ReceiverUsername string
	Amount           int
}
