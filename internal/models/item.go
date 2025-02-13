package models

type Item struct {
	Name  string `db:"name"`
	Price int    `db:"price"`
}

type PurchaseItem struct {
	Item  string `db:"item"`
	Count int    `db:"count"`
}
