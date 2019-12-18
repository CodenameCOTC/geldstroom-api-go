package transaction

import (
	"errors"
	"time"
)

type TransactionModel struct {
	Id          int       `json:"id"`
	Amount      int64     `json:"amount"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Type        string    `json:"type"`
	UserId      int       `json:"userId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type InsertDto struct {
	Amount      int64  `form:"amount"`
	Description string `form:"description"`
	Category    string `form:"category"`
	Type        string `form:"type"`
}

type updateDto struct {
	Amount      int64  `form:"amount"`
	Description string `form:"description"`
	Category    string `form:"category"`
	Type        string `form:"type"`
}

var (
	errInsert     = errors.New("Insert failure")
	errInsertCode = "TRANSACTION_001"
)
