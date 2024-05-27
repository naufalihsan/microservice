package common

import "syscall"

var (
	ApiGatewayService = "apiGateway"
	OrdersService     = "orders"
)

func EnvString(key, fallback string) string {
	if val, ok := syscall.Getenv(key); ok {
		return val
	}

	return fallback
}
