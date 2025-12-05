package middleware

import "github.com/marcofilho/go-ecommerce/src/internal/domain/entity"

type Permission string

const (
	// Product permissions
	PermissionCreateProduct Permission = "product:create"
	PermissionUpdateProduct Permission = "product:update"
	PermissionDeleteProduct Permission = "product:delete"
	PermissionViewProduct   Permission = "product:view"
	PermissionListProducts  Permission = "product:list"

	// Order permissions
	PermissionCreateOrder       Permission = "order:create"
	PermissionViewOrder         Permission = "order:view"
	PermissionListOrders        Permission = "order:list"
	PermissionUpdateOrderStatus Permission = "order:update_status"

	// Webhook permissions
	PermissionViewWebhookHistory Permission = "webhook:view_history"
)

var RolePermissions = map[entity.Role][]Permission{
	entity.RoleAdmin: {
		// Admins have all permissions
		PermissionCreateProduct,
		PermissionUpdateProduct,
		PermissionDeleteProduct,
		PermissionViewProduct,
		PermissionListProducts,
		PermissionCreateOrder,
		PermissionViewOrder,
		PermissionListOrders,
		PermissionUpdateOrderStatus,
		PermissionViewWebhookHistory,
	},
	entity.RoleCustomer: {
		// Customers can only view products and manage their own orders
		PermissionViewProduct,
		PermissionListProducts,
		PermissionCreateOrder,
		PermissionViewOrder,
		PermissionListOrders,
	},
}

func HasPermission(role entity.Role, permission Permission) bool {
	permissions, exists := RolePermissions[role]
	if !exists {
		return false
	}

	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}
