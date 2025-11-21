package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/marcofilho/go-ecommerce/src/internal/adapter/http/dto"
	"github.com/marcofilho/go-ecommerce/src/internal/domain/entity"
	"github.com/marcofilho/go-ecommerce/src/usecase/order"
)

type OrderHandler struct {
	useCase *order.UseCase
}

func NewOrderHandler(useCase *order.UseCase) *OrderHandler {
	return &OrderHandler{
		useCase: useCase,
	}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order with the provided products
// @Tags orders
// @Accept json
// @Produce json
// @Param order body dto.CreateOrderRequest true "Order information"
// @Success 201 {object} dto.OrderResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var products []order.CreateOrderItem
	for _, product := range req.Products {
		productID, err := uuid.Parse(product.ProductID)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid product ID")
			return
		}

		products = append(products, order.CreateOrderItem{
			ProductID: productID,
			Quantity:  product.Quantity,
		})
	}

	createdOrder, err := h.useCase.CreateOrder(r.Context(), req.CustomerID, products)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := dto.ToOrderResponse(createdOrder)
	respondJSON(w, http.StatusCreated, response)
}

// GetOrder godoc
// @Summary Get an order by ID
// @Description Get detailed information about a specific order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	order, err := h.useCase.GetOrder(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Order not found")
		return
	}

	response := dto.ToOrderResponse(order)

	respondJSON(w, http.StatusOK, response)
}

// ListOrders godoc
// @Summary List all orders
// @Description Get a paginated list of orders with optional filtering
// @Tags orders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(10)
// @Param status query string false "Filter by status (pending, cancelled, completed)"
// @Param payment_status query string false "Filter by payment status (unpaid, paid, failed)"
// @Success 200 {object} dto.OrderListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /orders [get]
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	statusStr := r.URL.Query().Get("status")
	paymentStatusStr := r.URL.Query().Get("payment_status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var status *entity.OrderStatus
	if statusStr != "" {
		s := entity.OrderStatus(statusStr)
		status = &s
	}

	var paymentStatus *entity.PaymentStatus
	if paymentStatusStr != "" {
		ps := entity.PaymentStatus(paymentStatusStr)
		paymentStatus = &ps
	}

	orders, total, err := h.useCase.ListOrders(r.Context(), page, pageSize, status, paymentStatus)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.ToOrderListResponse(orders, total, page, pageSize)

	respondJSON(w, http.StatusOK, response)
}

// UpdateOrderStatus godoc
// @Summary Update order status
// @Description Update the status of an existing order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param status body dto.UpdateOrderStatusRequest true "New status"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /orders/{id}/status [put]
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var req dto.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	newStatus := entity.OrderStatus(req.Status)
	order, err := h.useCase.UpdateOrderStatus(r.Context(), id, newStatus)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := dto.ToOrderResponse(order)

	respondJSON(w, http.StatusOK, response)
}
