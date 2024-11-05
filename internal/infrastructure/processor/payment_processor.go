package processor

import (
	"GO-API/internal/domain/model"
)

type PaymentProcessor struct {
}

func NewPaymentProcessor() *PaymentProcessor {
	return &PaymentProcessor{}
}

func (p *PaymentProcessor) Process(payment *model.Payment) error {
	return nil
}

func (p *PaymentProcessor) Cancel(payment *model.Payment) error {
	return nil
}
