package services

type GetBalanceParams struct {
	Token string
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

type GetTransactionsParams struct {
	Token string
}

type GetPurchasesParams struct {
	Token string
}

type BuyItemParams struct {
	Token string
	Item  string
}
