package common

import "syscall"

var (
	ApiGatewayService = "apiGateway"
	OrderService      = "orders"
	PaymentService    = "payments"
	InventoryService  = "inventory"
)

var (
	OrderStatusPending        = "pending"
	OrderStatusPaid           = "paid"
	OrderStatusWaitingPayment = "waiting_payment"
	OrderStatusDelivered      = "delivered"
)

func EnvString(key, fallback string) string {
	if val, ok := syscall.Getenv(key); ok {
		return val
	}

	return fallback
}
