package model

import "time"

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted  PaymentStatus = "completed"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusCanceled   PaymentStatus = "canceled"
)

type Payment struct {
	ID            string          `json:"id"`
	Amount        int64           `json:"amount"`
	Currency      string          `json:"currency"`
	Status        PaymentStatus   `json:"status"`
	Description   string          `json:"description"`
	CustomerID    string          `json:"customer_id"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	TransactionID string          `json:"transaction_id"`
	Metadata      PaymentMetadata `json:"metadata"`
}

type PaymentMetadata struct {
	OrderID       string `json:"order_id"`
	ProductID     string `json:"product_id"`
	PaymentMethod string `json:"payment_method"`
}

type PaymentRepository interface {
	Create(payment *Payment) error
	FindByID(id string) (*Payment, error)
	Update(payment *Payment) error
	List(limit int, offset int) ([]*Payment, error)
}

type PaymentProcessor interface {
	Process(payment *Payment) error
	Cancel(payment *Payment) error
}
