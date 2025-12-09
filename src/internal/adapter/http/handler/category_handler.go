package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/adapter/http/dto"
	"github.com/marcofilho/go-ecommerce/src/usecase/category"
)

type CategoryHandler struct {
	categoryService category.CategoryService
}

func NewCategoryHandler(categoryService category.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category (Admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Param category body dto.CategoryRequest true "Category details"
// @Success 201 {object} dto.CategoryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Security BearerAuth
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req dto.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	category, err := h.categoryService.CreateCategory(r.Context(), req.Name)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := dto.CategoryResponse{
		ID:   category.ID.String(),
		Name: category.Name,
	}

	respondJSON(w, http.StatusCreated, response)
}

// ListCategories godoc
// @Summary List all categories
// @Description Get a paginated list of categories
// @Tags categories
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.PaginatedResponse[dto.CategoryResponse]
// @Failure 500 {object} ErrorResponse
// @Router /categories [get]
func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	categories, total, err := h.categoryService.ListCategories(r.Context(), page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	categoryResponses := make([]dto.CategoryResponse, len(categories))
	for i, cat := range categories {
		categoryResponses[i] = dto.CategoryResponse{
			ID:   cat.ID.String(),
			Name: cat.Name,
		}
	}

	response := dto.PaginatedResponse[dto.CategoryResponse]{
		Data:     categoryResponses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	respondJSON(w, http.StatusOK, response)
}

// AssignCategoryToProduct godoc
// @Summary Assign category to product
// @Description Assign a category to a product (Admin only)
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body dto.AssignCategoryRequest true "Category assignment"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /products/{id}/categories [post]
func (h *CategoryHandler) AssignCategoryToProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req dto.AssignCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	if err := h.categoryService.AssignCategoryToProduct(r.Context(), productID, categoryID); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "Category assigned successfully"})
}

// RemoveCategoryFromProduct godoc
// @Summary Remove category from product
// @Description Remove a category from a product (Admin only)
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param category_id path string true "Category ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Security BearerAuth
// @Router /products/{id}/categories/{category_id} [delete]
func (h *CategoryHandler) RemoveCategoryFromProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	categoryIDStr := r.PathValue("category_id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	if err := h.categoryService.RemoveCategoryFromProduct(r.Context(), productID, categoryID); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "Category removed successfully"})
}

// GetProductCategories godoc
// @Summary Get product categories
// @Description Get all categories assigned to a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {array} dto.CategoryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/{id}/categories [get]
func (h *CategoryHandler) GetProductCategories(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	categories, err := h.categoryService.GetProductCategories(r.Context(), productID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	categoryResponses := make([]dto.CategoryResponse, len(categories))
	for i, cat := range categories {
		categoryResponses[i] = dto.CategoryResponse{
			ID:   cat.ID.String(),
			Name: cat.Name,
		}
	}

	respondJSON(w, http.StatusOK, categoryResponses)
}

type MessageResponse struct {
	Message string `json:"message"`
}
