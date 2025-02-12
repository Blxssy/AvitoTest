package models

import "time"

type Purchase struct {
	ID          int
	Username    string
	Item        string
	Price       int
	PurchasedAt time.Time
}
