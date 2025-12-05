package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/adapter/http/dto"
	productvariant "github.com/marcofilho/go-ecommerce/src/usecase/product_variant"
)

type ProductVariantHandler struct {
	useCase productvariant.ProductVariantService
}

func NewProductVariantHandler(useCase productvariant.ProductVariantService) *ProductVariantHandler {
	return &ProductVariantHandler{
		useCase: useCase,
	}
}

// CreateProductVariant godoc
// @Summary Create a new product variant
// @Description Create a new product variant with the provided information
// @Tags product_variants
// @Accept json
// @Produce json
// @Param product_variant body dto.ProductVariantRequest true "Product variant information"
// @Success 201 {object} dto.ProductVariantResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /product_variants [post]
func (h *ProductVariantHandler) CreateProductVariant(w http.ResponseWriter, r *http.Request) {
	var req dto.ProductVariantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	productVariant, err := h.useCase.CreateProductVariant(r.Context(), productID, req.VariantName, req.VariantValue, req.PriceOverride, req.Quantity)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := dto.ToProductVariantResponse(productVariant)
	respondJSON(w, http.StatusCreated, response)
}

// GetProductVariant godoc
// @Summary Get a product variant by ID
// @Description Get detailed information about a specific product variant
// @Tags product_variants
// @Accept json
// @Produce json
// @Param id path string true "Product Variant ID"
// @Success 200 {object} dto.ProductVariantResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /product_variants/{id} [get]
func (h *ProductVariantHandler) GetProductVariant(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product variant ID")
		return
	}

	productVariant, err := h.useCase.GetProductVariant(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Product variant not found")
		return
	}

	response := dto.ToProductVariantResponse(productVariant)
	respondJSON(w, http.StatusOK, response)
}

// ListProductVariants godoc
// @Summary List all product variants for a product
// @Description Get a paginated list of product variants for a specific product
// @Tags product_variants
// @Accept json
// @Produce json
// @Param product_id query string true "Product ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(10)
// @Success 200 {object} dto.ProductVariantListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /product_variants [get]
func (h *ProductVariantHandler) ListProductVariants(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.URL.Query().Get("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	variants, total, err := h.useCase.ListProductVariants(r.Context(), productID, page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.ToProductVariantListResponse(variants, total, page, pageSize)
	respondJSON(w, http.StatusOK, response)
}

// UpdateProductVariant godoc
// @Summary Update a product variant
// @Description Update an existing product variant's information
// @Tags product_variants
// @Accept json
// @Produce json
// @Param id path string true "Product Variant ID"
// @Param product_variant body dto.ProductVariantRequest true "Product Variant information"
// @Success 200 {object} dto.ProductVariantResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /product_variants/{id} [put]
func (h *ProductVariantHandler) UpdateProductVariant(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product variant ID")
		return
	}

	var req dto.ProductVariantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	productVariant, err := h.useCase.UpdateProductVariant(r.Context(), id, req.VariantName, req.VariantValue, req.PriceOverride, req.Quantity)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := dto.ToProductVariantResponse(productVariant)
	respondJSON(w, http.StatusOK, response)
}

// DeleteProductVariant godoc
// @Summary Delete a product variant
// @Description Delete a product variant by ID
// @Tags product_variants
// @Accept json
// @Produce json
// @Param id path string true "Product Variant ID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /product_variants/{id} [delete]
func (h *ProductVariantHandler) DeleteProductVariant(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product variant ID")
		return
	}

	if err := h.useCase.DeleteProductVariant(r.Context(), id); err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
