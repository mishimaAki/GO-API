package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"GO-API/internal/pkg/logger"
	"GO-API/internal/usecase"
)

type PaymentHandler struct {
	paymentUseCase *usecase.PaymentUseCase
}

func NewPaymentHandler(pu *usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{
		paymentUseCase: pu,
	}
}

type CreatePaymentRequest struct {
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	Description   string `json:"description"`
	CustomerID    string `json:"customer_id"`
	PaymentMethod string `jsong:"payment_method"`
	OrderID       string `json:"order_id"`
}

func (h *PaymentHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/payments", h.CreatePayment).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/payments/{id}", h.GetPayment).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/payments", h.ListPayments).Methods(http.MethodGet)
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received create payment request")

	var req CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode request body: %v", err)
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	input := usecase.CreatePaymentInput{
		Amount:        req.Amount,
		Currency:      req.Currency,
		Description:   req.Description,
		CustomerID:    req.CustomerID,
		PaymentMethod: req.PaymentMethod,
		OrderID:       req.OrderID,
	}

	payment, err := h.paymentUseCase.CreatePayment(r.Context(), input)
	if err != nil {
		logger.Error("Failed to create payment: %v", err)
		handleError(w, err)
		return
	}

	logger.Info("Successfully created payment: ID=%s", payment.ID)
	writeJSON(w, http.StatusCreated, payment)
}

func (h *PaymentHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	logger.Info("GetPayment handler called")

	vars := mux.Vars(r)
	id := vars["id"]
	logger.Info("GetPayment handler called")

	payment, err := h.paymentUseCase.GetPayment(r.Context(), id)
	if err != nil {
		logger.Error("Error getting payment; %v", err)
		handleError(w, err)
		return
	}

	logger.Info("Successfully retrieved payment: %+v", payment)
	writeJSON(w, http.StatusOK, payment)
}

func (h *PaymentHandler) ListPayments(w http.ResponseWriter, r *http.Request) {
	logger.Info("ListPayments handler called")

	limit := 10
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = 1
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	logger.Info("Fetching payments with limit=%d, offset=%d", limit, offset)

	payments, err := h.paymentUseCase.ListPayments(r.Context(), limit, offset)
	if err != nil {
		logger.Error("Failed to fetch payments: %v", err)
		handleError(w, err)
		return
	}

	logger.Info("Successfully fetched %d payments", len(payments))
	writeJSON(w, http.StatusOK, payments)
}
