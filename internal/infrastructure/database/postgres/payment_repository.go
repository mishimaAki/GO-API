package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"GO-API/internal/domain/model"
	"GO-API/internal/pkg/logger"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

const createTableSQL = `
CREATE TABLE IF NOT EXISTS payments (
	id TEXT PRIMARY KEY,
	amount BIGINT NOT NULL,
	currency TEXT NOT NULL,
	status TEXT NOT NULL,
	description TEXT,
	customer_id TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	transaction_id TEXT NOT NULL UNIQUE,
	metadata JSONB NOT NULL);`

func (r *PaymentRepository) InitTable() error {
	_, err := r.db.Exec(createTableSQL)
	return err
}

func (r *PaymentRepository) Create(payment *model.Payment) error {
	logger.Info("Creating payment: ID=%s, Amount=%d, Currency=%s",
		payment.ID, payment.Amount, payment.Currency)

	query := `
		INSERT INTO payments (
			id, amount, currency, status, description, customer_id,
			created_at, updated_at, transaction_id, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	metadataJSON, err := json.Marshal(payment.Metadata)
	if err != nil {
		logger.Error("Failed to marshal metadata: %v", err)
		return err
	}
	_, err = r.db.Exec(
		query,
		payment.ID,
		payment.Amount,
		payment.Currency,
		payment.Status,
		payment.Description,
		payment.CustomerID,
		payment.CreatedAt,
		payment.UpdatedAt,
		payment.TransactionID,
		metadataJSON,
	)

	if err != nil {
		logger.Error("Failed to execute insert query: %v", err)
		return fmt.Errorf("error creating payment; %w", err)
	}

	logger.Info("Successfully created payment: ID=%s", payment.ID)
	return nil
}

func (r *PaymentRepository) FindByID(id string) (*model.Payment, error) {
	logger.Info("Executing FindByID query for ID: %s", id)

	query := `
		SELECT id, amount, currency, status, description, customer_id, 
				created_at, updated_at, transaction_id, metadata
		FROM payments
		WHERE id = $1`

	var payment model.Payment
	var metadataBytes []byte

	err := r.db.QueryRow(query, id).Scan(
		&payment.ID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.Description,
		&payment.CustomerID,
		&payment.CreatedAt,
		&payment.UpdatedAt,
		&payment.TransactionID,
		&metadataBytes,
	)

	if err == sql.ErrNoRows {
		logger.Error("Payment not found; %s", id)
		return nil, model.NewNotFoundError(("payment not found"))
	}

	if err != nil {
		logger.Error("Database error; %v", err)
		return nil, fmt.Errorf("error finding payment: %w", err)
	}

	if err := json.Unmarshal(metadataBytes, &payment.Metadata); err != nil {
		logger.Error("Failed to unmarshal metadata; %v", err)
		return nil, err
	}

	logger.Debug("Successfully found payment: %+v", payment)
	return &payment, nil
}

func (r *PaymentRepository) List(limit int, offset int) ([]*model.Payment, error) {
	logger.Info("Executing List query with limit=%d, offset=%d", limit, offset)

	query := `
		SELECT id, amount, currency, status, description, customer_id,
			  created_at, updated_at, transaction_id, metadata
		FROM payments
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		logger.Error("Failed to execute list query: %v", err)
		return nil, fmt.Errorf("error listing payments: %w", err)
	}
	defer rows.Close()

	var payments []*model.Payment
	for rows.Next() {
		var payment model.Payment
		var metadataBytes []byte
		err := rows.Scan(
			&payment.ID,
			&payment.Amount,
			&payment.Currency,
			&payment.Status,
			&payment.Description,
			&payment.CustomerID,
			&payment.CreatedAt,
			&payment.UpdatedAt,
			&payment.TransactionID,
			&metadataBytes,
		)

		if err != nil {
			logger.Error("Failed to scan payment row: %v", err)
			return nil, fmt.Errorf("error scanning payment row: %w", err)
		}

		if err := json.Unmarshal(metadataBytes, &payment.Metadata); err != nil {
			logger.Error("Failed to unmarshal metadata: %v", err)
			return nil, fmt.Errorf("error unmarshaling metadata: %w", err)
		}

		payments = append(payments, &payment)
	}

	logger.Info("Successfully retrieved %d payments", len(payments))
	return payments, nil
}

func (r *PaymentRepository) Update(payment *model.Payment) error {
	logger.Info("Updating payment: ID=%s", payment.ID)

	query := `
		UPDATE payments
		SET amount = $1,
			currency = $2,
			status = $3,
			description = $4,
			customer_id = $5,
			updated_at = $6,
			transaction_id = $7,
			metadata = $8
		WHERE id = $9`

	result, err := r.db.Exec(
		query,
		payment.Amount,
		payment.Currency,
		payment.Status,
		payment.Description,
		payment.CustomerID,
		time.Now(),
		payment.TransactionID,
		payment.Metadata,
		payment.ID,
	)

	if err != nil {
		logger.Error("Failed to execute update query: %v", err)
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.Error("Failed to get affected rows: %v", err)
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rows == 0 {
		logger.Error("Payment not found for update: %s", payment.ID)
		return model.NewNotFoundError("payment not found")
	}

	logger.Info("Successfully updated payment: ID=%s", payment.ID)
	return nil
}
