package repo

type GetBalanceParams struct {
	UserID uint32
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
