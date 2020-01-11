package entity

import "time"

type Transaction struct {
	Id          string    `json:"id"`
	Amount      int64     `json:"amount"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Type        string    `json:"type"`
	UserId      string    `json:"userId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type TotalAmount struct {
	Income  int `json:"income"`
	Expense int `json:"expense"`
}
