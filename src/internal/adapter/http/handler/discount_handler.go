package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/adapter/http/dto"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/usecase/discount"
)

type DiscountHandler struct {
	useCase discount.DiscountService
}

func NewDiscountHandler(useCase discount.DiscountService) *DiscountHandler {
	return &DiscountHandler{
		useCase: useCase,
	}
}

// CreateDiscount godoc
// @Summary Create a new discount
// @Description Create a new discount/promo code (Admin only)
// @Tags discounts
// @Accept json
// @Produce json
// @Param discount body dto.DiscountRequest true "Discount information"
// @Success 201 {object} dto.DiscountResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /discounts [post]
func (h *DiscountHandler) CreateDiscount(w http.ResponseWriter, r *http.Request) {
	var req dto.DiscountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	discountType := entity.DiscountType(req.DiscountType)
	if discountType != entity.Percentage && discountType != entity.Amount {
		respondError(w, http.StatusBadRequest, "Invalid discount_type. Must be 'percentage' or 'amount'")
		return
	}

	if discountType == entity.Percentage && (req.Value < 0 || req.Value > 100) {
		respondError(w, http.StatusBadRequest, "Percentage discount must be between 0 and 100")
		return
	}
	if req.Value <= 0 {
		respondError(w, http.StatusBadRequest, "Discount value must be greater than 0")
		return
	}

	discount := &entity.Discount{
		PromoCode:    req.PromoCode,
		DiscountType: discountType,
		Value:        req.Value,
		Active:       req.Active,
	}

	if err := h.useCase.CreateDiscount(r.Context(), discount); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := dto.ToDiscountResponse(discount)
	respondJSON(w, http.StatusCreated, response)
}

// GetDiscount godoc
// @Summary Get a discount by ID
// @Description Get detailed information about a specific discount (Admin only)
// @Tags discounts
// @Accept json
// @Produce json
// @Param id path string true "Discount ID"
// @Success 200 {object} dto.DiscountResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /discounts/{id} [get]
func (h *DiscountHandler) GetDiscount(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid discount ID")
		return
	}

	discount, err := h.useCase.GetDiscountByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Discount not found")
		return
	}

	response := dto.ToDiscountResponse(discount)
	respondJSON(w, http.StatusOK, response)
}

// UpdateDiscount godoc
// @Summary Update a discount
// @Description Update an existing discount/promo code (Admin only)
// @Tags discounts
// @Accept json
// @Produce json
// @Param id path string true "Discount ID"
// @Param discount body dto.DiscountRequest true "Updated discount information"
// @Success 200 {object} dto.DiscountResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /discounts/{id} [put]
func (h *DiscountHandler) UpdateDiscount(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid discount ID")
		return
	}

	existingDiscount, err := h.useCase.GetDiscountByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Discount not found")
		return
	}

	var req dto.DiscountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	discountType := entity.DiscountType(req.DiscountType)
	if discountType != entity.Percentage && discountType != entity.Amount {
		respondError(w, http.StatusBadRequest, "Invalid discount_type. Must be 'percentage' or 'amount'")
		return
	}

	if discountType == entity.Percentage && (req.Value < 0 || req.Value > 100) {
		respondError(w, http.StatusBadRequest, "Percentage discount must be between 0 and 100")
		return
	}
	if req.Value <= 0 {
		respondError(w, http.StatusBadRequest, "Discount value must be greater than 0")
		return
	}

	existingDiscount.PromoCode = req.PromoCode
	existingDiscount.DiscountType = discountType
	existingDiscount.Value = req.Value
	existingDiscount.Active = req.Active

	if err := h.useCase.UpdateDiscount(r.Context(), existingDiscount); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := dto.ToDiscountResponse(existingDiscount)
	respondJSON(w, http.StatusOK, response)
}

// ValidatePromoCode godoc
// @Summary Validate a promo code
// @Description Validate a promo code and return discount information if active
// @Tags discounts
// @Accept json
// @Produce json
// @Param request body dto.ApplyDiscountRequest true "Promo code to validate"
// @Success 200 {object} dto.DiscountResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /discounts/validate [post]
func (h *DiscountHandler) ValidatePromoCode(w http.ResponseWriter, r *http.Request) {
	var req dto.ApplyDiscountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.PromoCode == "" {
		respondError(w, http.StatusBadRequest, "Promo code is required")
		return
	}

	discount, err := h.useCase.GetDiscountByPromoCode(r.Context(), req.PromoCode)
	if err != nil {
		respondError(w, http.StatusNotFound, "Invalid or inactive promo code")
		return
	}

	response := dto.ToDiscountResponse(discount)
	respondJSON(w, http.StatusOK, response)
}
