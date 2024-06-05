package common

import "syscall"

var (
	ApiGatewayService = "apiGateway"
	OrderService      = "orders"
	PaymentService    = "payments"
)

var (
	OrderStatusPending        = "pending"
	OrderStatusPaid           = "paid"
	OrderStatusWaitingPayment = "waiting_payment"
)

func EnvString(key, fallback string) string {
	if val, ok := syscall.Getenv(key); ok {
		return val
	}

	return fallback
}
