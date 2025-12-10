package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/usecase/payment"
)

type PaymentHandler struct {
	paymentUC     payment.PaymentService
	webhookSecret string
}

func NewPaymentHandler(paymentUC payment.PaymentService, webhookSecret string) *PaymentHandler {
	return &PaymentHandler{
		paymentUC:     paymentUC,
		webhookSecret: webhookSecret,
	}
}

// PaymentWebhookHandler handles incoming payment webhooks
// @Summary Process payment webhook
// @Description Receives payment status updates from payment processor with HMAC signature verification and replay attack prevention
// @Tags payments
// @Accept json
// @Produce json
// @Param X-Payment-Signature header string true "HMAC-SHA256 signature of the request body"
// @Param webhook body entity.PaymentWebhookRequest true "Payment webhook data with timestamp"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Unauthorized - Invalid signature or timestamp"
// @Router /payment-webhook [post]
func (h *PaymentHandler) PaymentWebhookHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}
	defer r.Body.Close()

	signature := r.Header.Get("X-Payment-Signature")
	if signature == "" {
		respondError(w, http.StatusUnauthorized, "Missing payment signature")
		return
	}

	if !h.verifySignature(body, signature) {
		respondError(w, http.StatusUnauthorized, "Invalid payment signature")
		return
	}

	var req entity.PaymentWebhookRequest
	if err := json.Unmarshal(body, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if !h.verifyTimestamp(req.Timestamp) {
		respondError(w, http.StatusUnauthorized, "Request timestamp is too old or invalid")
		return
	}

	if err := h.paymentUC.ProcessWebhook(r.Context(), &req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Payment webhook processed successfully",
	})
}

// GetWebhookHistoryHandler retrieves webhook history for an order
// @Summary Get payment webhook history
// @Description Retrieves all payment webhook events for a specific order
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {array} entity.WebhookLog
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id}/payment-history [get]
func (h *PaymentHandler) GetWebhookHistoryHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if _, err := uuid.Parse(idStr); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	logs, err := h.paymentUC.GetWebhookHistory(r.Context(), idStr)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, logs)
}

// verifySignature validates the HMAC signature of the webhook payload
func (h *PaymentHandler) verifySignature(payload []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(h.webhookSecret))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func (h *PaymentHandler) verifyTimestamp(timestamp int64) bool {
	if timestamp == 0 {
		return false
	}

	webhookTime := time.Unix(timestamp, 0)
	now := time.Now()

	if webhookTime.After(now.Add(5 * time.Minute)) {
		return false
	}

	if webhookTime.Before(now.Add(-5 * time.Minute)) {
		return false
	}

	return true
}
