package dto

import (
	"mime/multipart"
	"time"

	"github.com/shopspring/decimal"
)

type ExpenseStatus string

const (
	ExpenseStatusPaid    ExpenseStatus = "paid"
	ExpenseStatusPending ExpenseStatus = "pending"
)

type ExpenseCategoryDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ExpenseAttachmentDTO struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

type ExpenseDTO struct {
	ID          int64                  `json:"id"`
	Description string                 `json:"description"`
	CategoryID  int64                  `json:"category_id"`
	Amount      decimal.Decimal        `json:"amount"`
	Date        time.Time              `json:"date"`
	Status      ExpenseStatus          `json:"status"`
	Category    *ExpenseCategoryDTO    `json:"category,omitempty"`
	Attachments []ExpenseAttachmentDTO `json:"attachments,omitempty"`
}

type ExpenseCreateUpdateDTO struct {
	Description string          `json:"description"`
	CategoryID  int64           `json:"category_id"`
	Amount      decimal.Decimal `json:"amount"`
	Date        time.Time       `json:"date"`
	Status      ExpenseStatus   `json:"status"`

	Attachments []ExpenseAttachmentDTO  `json:"attachments,omitempty"`
	Files       []*multipart.FileHeader `json:"files,omitempty"`
}
