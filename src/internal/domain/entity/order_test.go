package entity

import (
	"testing"

	"github.com/google/uuid"
)

func TestOrder_CalculateTotal(t *testing.T) {
	tests := []struct {
		name  string
		order Order
		want  float64
	}{
		{
			name: "single item",
			order: Order{
				Products: []OrderItem{
					{Price: 100.00, Quantity: 2},
				},
			},
			want: 200.00,
		},
		{
			name: "multiple items",
			order: Order{
				Products: []OrderItem{
					{Price: 100.00, Quantity: 2},
					{Price: 50.00, Quantity: 3},
				},
			},
			want: 350.00,
		},
		{
			name: "empty order",
			order: Order{
				Products: []OrderItem{},
			},
			want: 0.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.order.CalculateTotal()
			if tt.order.TotalPrice != tt.want {
				t.Errorf("CalculateTotal() = %v, want %v", tt.order.TotalPrice, tt.want)
			}
		})
	}
}

func TestOrder_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name      string
		current   OrderStatus
		newStatus OrderStatus
		want      bool
	}{
		{
			name:      "pending to completed",
			current:   Pending,
			newStatus: Completed,
			want:      true,
		},
		{
			name:      "pending to cancelled",
			current:   Pending,
			newStatus: Cancelled,
			want:      true,
		},
		{
			name:      "completed to cancelled",
			current:   Completed,
			newStatus: Cancelled,
			want:      false,
		},
		{
			name:      "cancelled to completed",
			current:   Cancelled,
			newStatus: Completed,
			want:      false,
		},
		{
			name:      "same status",
			current:   Pending,
			newStatus: Pending,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := Order{Status: tt.current}
			err := order.CanTransitionTo(tt.newStatus)
			got := err == nil
			if got != tt.want {
				t.Errorf("CanTransitionTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrder_UpdateStatus(t *testing.T) {
	tests := []struct {
		name      string
		current   OrderStatus
		newStatus OrderStatus
		wantErr   bool
	}{
		{
			name:      "valid transition",
			current:   Pending,
			newStatus: Completed,
			wantErr:   false,
		},
		{
			name:      "invalid transition",
			current:   Completed,
			newStatus: Cancelled,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := &Order{Status: tt.current}
			err := order.UpdateStatus(tt.newStatus)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && order.Status != tt.newStatus {
				t.Errorf("UpdateStatus() status = %v, want %v", order.Status, tt.newStatus)
			}
		})
	}
}

func TestOrderItem_Validate(t *testing.T) {
	validProductID := uuid.New()

	tests := []struct {
		name    string
		item    OrderItem
		wantErr bool
	}{
		{
			name: "valid item",
			item: OrderItem{
				ProductID: validProductID,
				Quantity:  2,
				Price:     100.00,
			},
			wantErr: false,
		},
		{
			name: "zero quantity",
			item: OrderItem{
				ProductID: validProductID,
				Quantity:  0,
				Price:     100.00,
			},
			wantErr: true,
		},
		{
			name: "negative quantity",
			item: OrderItem{
				ProductID: validProductID,
				Quantity:  -1,
				Price:     100.00,
			},
			wantErr: true,
		},
		{
			name: "negative price",
			item: OrderItem{
				ProductID: validProductID,
				Quantity:  2,
				Price:     -100.00,
			},
			wantErr: true,
		},
		{
			name: "nil product ID",
			item: OrderItem{
				ProductID: uuid.Nil,
				Quantity:  2,
				Price:     100.00,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrderItem_Subtotal(t *testing.T) {
	tests := []struct {
		name string
		item OrderItem
		want float64
	}{
		{
			name: "calculate subtotal",
			item: OrderItem{
				Quantity: 3,
				Price:    100.00,
			},
			want: 300.00,
		},
		{
			name: "zero quantity",
			item: OrderItem{
				Quantity: 0,
				Price:    99.99,
			},
			want: 0.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.item.Subtotal(); got != tt.want {
				t.Errorf("Subtotal() = %v, want %v", got, tt.want)
			}
		})
	}
}
