package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/adapter/http/dto"
	"github.com/marcofilho/go-ecommerce/src/usecase/product"
)

type ProductHandler struct {
	useCase *product.UseCase
}

func NewProductHandler(useCase *product.UseCase) *ProductHandler {
	return &ProductHandler{
		useCase: useCase,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req dto.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	product, err := h.useCase.CreateProduct(r.Context(), req.Name, req.Description, req.Price, req.Quantity)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := dto.ToProductResponse(product)
	respondJSON(w, http.StatusCreated, response)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.useCase.GetProduct(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Product not found")
		return
	}

	response := dto.ToProductResponse(product)
	respondJSON(w, http.StatusOK, response)
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	inStockOnly := r.URL.Query().Get("in_stock_only") == "true"

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	products, total, err := h.useCase.ListProducts(r.Context(), page, pageSize, inStockOnly)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.ToProductListResponse(products, total, page, pageSize)
	respondJSON(w, http.StatusOK, response)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req dto.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	product, err := h.useCase.UpdateProduct(r.Context(), id, req.Name, req.Description, req.Price, req.Quantity)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := dto.ToProductResponse(product)
	respondJSON(w, http.StatusOK, response)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	if err := h.useCase.DeleteProduct(r.Context(), id); err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
