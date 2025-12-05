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
	useCase product.ProductService
}

func NewProductHandler(useCase product.ProductService) *ProductHandler {
	return &ProductHandler{
		useCase: useCase,
	}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the provided information
// @Tags products
// @Accept json
// @Produce json
// @Param product body dto.ProductRequest true "Product information"
// @Success 201 {object} dto.ProductResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /products [post]
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

// GetProduct godoc
// @Summary Get a product by ID
// @Description Get detailed information about a specific product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /products/{id} [get]
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

// ListProducts godoc
// @Summary List all products
// @Description Get a paginated list of products with optional filtering
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(10)
// @Param in_stock_only query bool false "Filter products in stock only" default(true)
// @Success 200 {object} dto.ProductListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /products [get]
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	inStockOnlyParam := r.URL.Query().Get("in_stock_only")
	inStockOnly := true
	if inStockOnlyParam == "false" {
		inStockOnly = false
	}

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

// UpdateProduct godoc
// @Summary Update a product
// @Description Update an existing product's information
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body dto.ProductRequest true "Product information"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /products/{id} [put]
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

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /products/{id} [delete]
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
