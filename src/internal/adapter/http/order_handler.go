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

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest
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

	response := dto.ToOrderResponse(createdOrder)
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

	response := dto.ToOrderResponse(order)
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

	response := dto.ToOrderListResponse(orders, total, page, pageSize)
	respondJSON(w, http.StatusOK, response)
}

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
