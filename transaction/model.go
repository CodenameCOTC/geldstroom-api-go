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

type insertDto struct {
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
	errInsert                = errors.New("Insert failure")
	errInsertCode            = "TRANSACTION_0001"
	errInvalidDateRange      = errors.New("Invalid date range, must be one of DAILY | WEEKLY | MONTHLY")
	errInvalidDateRangeCode  = "TRANSASCTION_0002"
	errInvalidDateFormat     = errors.New("Invalid date format | date format must be YYYY-MM-DD")
	errInvalidDateFormatCode = "TRANSACTION_0003"
)

var (
	oneWeekRange  = "WEEKLY"
	oneDayRange   = "DAILY"
	oneMonthRange = "MONTHLY"
)
