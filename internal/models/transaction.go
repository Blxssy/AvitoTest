package models

import "time"

type Transaction struct {
	ID               uint
	SenderUsername   string
	ReceiverUsername string
	Amount           int
	CreatedAt        time.Time
}
