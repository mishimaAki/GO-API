package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"GO-API/internal/domain/model"
)

type TransactionIDGenerator interface {
	Generate() (string, error)
}

type DefaultTransactionIDGenerator struct {
	prefix string
}

func NewTransactionIDGenerator(prefix string) *DefaultTransactionIDGenerator {
	return &DefaultTransactionIDGenerator{
		prefix: prefix,
	}
}

func (g *DefaultTransactionIDGenerator) Generate() (string, error) {
	timestamp := time.Now().Format("20060102150405")

	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	randomHex := hex.EncodeToString(randomBytes)
	transactionID := fmt.Sprintf("%s_%s_%s", g.prefix, timestamp, randomHex)

	return transactionID, nil
}

type PaymentTransactionIDGenerator struct {
	generator TransactionIDGenerator
}

func NewPaymentTransactionIDGenerator() *PaymentTransactionIDGenerator {
	return &PaymentTransactionIDGenerator{
		generator: NewTransactionIDGenerator("PAY"),
	}
}

func (g *PaymentTransactionIDGenerator) Generate() (string, error) {
	return g.generator.Generate()
}

func ValidateTransactionID(transactionID string) error {
	if len(transactionID) < 25 {
		return model.NewValidationError("invalid transaction id format")
	}

	if transactionID[:3] != "PAY" {
		return model.NewValidationError("invalid transaction id prefix")
	}

	return nil
}
