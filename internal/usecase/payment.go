package usecase

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"GO-API/internal/domain/model"
	"GO-API/internal/domain/service"
	"GO-API/internal/gateway"
	"GO-API/internal/pkg/logger"
)

type PaymentUseCase struct {
	repo          gateway.PaymentRepository
	processor     gateway.PaymentProcessor
	txIDGenerator *service.PaymentTransactionIDGenerator
}

const (
	CurrencyJPY = "JPY"
	CurrencyUSD = "USD"
)

const (
	MaxDescriptionLength = 500
	MaxCustomerIDLength  = 100
	MinAmount            = 1
	MaxAmount            = 10000000
)

func NewPaymentUseCase(repo gateway.PaymentRepository, processor gateway.PaymentProcessor) *PaymentUseCase {
	return &PaymentUseCase{
		repo:          repo,
		processor:     processor,
		txIDGenerator: service.NewPaymentTransactionIDGenerator(),
	}
}

type CreatePaymentInput struct {
	Amount        int64
	Currency      string
	Description   string
	CustomerID    string
	PaymentMethod string
	OrderID       string
}

func (uc *PaymentUseCase) CreatePayment(ctx context.Context, input CreatePaymentInput) (*model.Payment, error) {
	logger.Info("Creating payment with amount=%d currency=%s", input.Amount, input.Currency)

	if err := validateCreatePaymentInput(input); err != nil {
		logger.Error("Payment validation failed: %v", err)
		log.Printf("Validation error: %v", err)
		return nil, err
	}

	transactionID, err := uc.txIDGenerator.Generate()
	if err != nil {
		logger.Error("Failed to generate transaction ID: %v", err)
		log.Printf("Transaction ID generation error: %v", err)
		return nil, model.NewInternalError(err)
	}
	log.Printf("Generated transaction ID: %s", transactionID)

	payment := &model.Payment{
		ID:            uuid.New().String(),
		Amount:        input.Amount,
		Currency:      input.Currency,
		Status:        model.PaymentStatusPending,
		Description:   input.Description,
		CustomerID:    input.CustomerID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		TransactionID: transactionID,
		Metadata: model.PaymentMetadata{
			OrderID:       input.OrderID,
			PaymentMethod: input.PaymentMethod,
		},
	}
	log.Printf("Created payment object: %+v", payment)

	if err := uc.repo.Create(payment); err != nil {
		logger.Error("Database error: %v", err)
		return nil, model.NewInternalError(err)
	}

	if err := uc.processor.Process(payment); err != nil {
		logger.Error("Processing error: %v", err)
		return nil, model.NewInternalError(err)
	}

	logger.Info("Successfully processed payment: ID=%s", payment.ID)
	return payment, nil
}

func validateCreatePaymentInput(input CreatePaymentInput) error {
	if input.Amount <= 0 {
		return model.NewValidationError("amount must be positive")
	}
	if input.Currency == "" {
		return model.NewValidationError("currency is required")
	}
	if input.CustomerID == "" {
		return model.NewValidationError("customer_id is required")
	}

	if input.Amount <= 0 {
		logger.Error("Invalid amount: %d", input.Amount)
		return model.NewValidationError("amount must be positive")
	}

	if input.Amount > MaxAmount {
		logger.Error("Amount exceeds maximum: %d", input.Amount)
		return model.NewValidationError("amount exceeds maximum allowed")
	}

	if input.Currency == "" {
		logger.Error("Currency is empty")
		return model.NewValidationError("currency is required")
	}

	validCurrencies := map[string]bool{
		CurrencyJPY: true,
		CurrencyUSD: true,
	}
	if !validCurrencies[input.Currency] {
		logger.Error("Unsupported currency: %s", input.Currency)
		return model.NewValidationError("unsupported currency")
	}

	if err := validatePaymentMethod(input.PaymentMethod); err != nil {
		logger.Error("Invalid payment method: %s", input.PaymentMethod)
		return err
	}

	logger.Debug("Validation passed for payment: amount=%d currency=%s customer=%s",
		input.Amount, input.Currency, input.CustomerID)

	return nil
}

func validatePaymentMethod(method string) error {
	validMethods := map[string]bool{
		"credit_card":       true,
		"bank_transfer":     true,
		"convenience_store": true,
	}

	if method == "" {
		return model.NewValidationError("payment_method is required")
	}

	if !validMethods[method] {
		return model.NewValidationError("unsupported payment method")
	}

	return nil
}

func (uc *PaymentUseCase) GetPayment(ctx context.Context, id string) (*model.Payment, error) {
	logger.Info("Getting payment by ID; %s", id)

	payment, err := uc.repo.FindByID(id)
	if err != nil {
		logger.Error("Faild to find payment; %v", err)
		return nil, err
	}

	logger.Debug("Found payment; %+v", payment)
	return payment, nil
}

func (uc *PaymentUseCase) ListPayments(ctx context.Context, limit, offset int) ([]*model.Payment, error) {
	logger.Info("Listing payments with limit=%d offset=%d", limit, offset)

	if limit <= 0 {
		limit = 10
		logger.Debug("Using default limit value: 10")
	}

	if offset < 0 {
		offset = 0
		logger.Debug("Adjusting negative offset to 0")
	}

	payments, err := uc.repo.List(limit, offset)
	if err != nil {
		logger.Error("Failed to list payments: %v", err)
		return nil, err
	}

	logger.Info("Successfully retrieved %d payments", len(payments))
	return payments, nil
}
