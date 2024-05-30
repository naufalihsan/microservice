package common

import "syscall"

var (
	ApiGatewayService = "apiGateway"
	OrderService      = "orders"
	PaymentService    = "payments"
)

func EnvString(key, fallback string) string {
	if val, ok := syscall.Getenv(key); ok {
		return val
	}

	return fallback
}
