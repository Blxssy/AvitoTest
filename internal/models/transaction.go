package models

import "time"

type Transaction struct {
	ID               uint32
	SenderUsername   string
	ReceiverUsername string
	Amount           int
	CreatedAt        time.Time
}
