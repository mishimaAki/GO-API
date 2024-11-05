package gateway

import (
	"GO-API/internal/domain/model"
)

type PaymentRepository interface {
	Create(payment *model.Payment) error
	FindByID(id string) (*model.Payment, error)
	Update(payment *model.Payment) error
	List(limit int, offset int) ([]*model.Payment, error)
}

type PaymentProcessor interface {
	Process(payment *model.Payment) error
	Cancel(payment *model.Payment) error
}
