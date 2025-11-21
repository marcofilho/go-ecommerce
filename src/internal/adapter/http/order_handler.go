package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
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

type CreateOrderRequest struct {
	CustomerID int                      `json:"customer_id"`
	Products   []CreateOrderItemRequest `json:"products"`
}

type CreateOrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}

type OrderItemResponse struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Subtotal  float64 `json:"subtotal"`
}

type OrderResponse struct {
	ID            string              `json:"id"`
	CustomerID    int                 `json:"customer_id"`
	Items         []OrderItemResponse `json:"items"`
	TotalPrice    float64             `json:"total_price"`
	Status        string              `json:"status"`
	PaymentStatus string              `json:"payment_status"`
	CreatedAt     string              `json:"created_at"`
	UpdatedAt     string              `json:"updated_at"`
}

type OrderListResponse struct {
	Orders   []OrderResponse `json:"orders"`
	Total    int             `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var items []order.CreateOrderItem
	for _, item := range req.Products {
		productID, err := uuid.Parse(item.ProductID)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid product ID")
			return
		}

		items = append(items, order.CreateOrderItem{
			ProductID: productID,
			Quantity:  item.Quantity,
		})
	}

	createdOrder, err := h.useCase.CreateOrder(r.Context(), req.CustomerID, items)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := orderToResponse(createdOrder)
	respondJSON(w, http.StatusCreated, response)
}

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

	response := orderToResponse(order)
	respondJSON(w, http.StatusOK, response)
}

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

	var orderResponses []OrderResponse
	for _, ord := range orders {
		orderResponses = append(orderResponses, orderToResponse(ord))
	}

	response := OrderListResponse{
		Orders:   orderResponses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var req UpdateOrderStatusRequest
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

	response := orderToResponse(order)
	respondJSON(w, http.StatusOK, response)
}

func orderToResponse(ord *entity.Order) OrderResponse {
	var items []OrderItemResponse
	for _, item := range ord.Items {
		items = append(items, OrderItemResponse{
			ProductID: item.ProductID.String(),
			Quantity:  item.Quantity,
			Price:     item.Price,
			Subtotal:  item.Subtotal(),
		})
	}

	return OrderResponse{
		ID:            ord.ID.String(),
		CustomerID:    ord.CustomerID,
		Items:         items,
		TotalPrice:    ord.TotalPrice,
		Status:        string(ord.Status),
		PaymentStatus: string(ord.PaymentStatus),
		CreatedAt:     ord.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     ord.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
