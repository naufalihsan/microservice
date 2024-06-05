package gateway

import (
	"context"
)

type OrderGateaway interface {
	UpdateOrderAfterPaymentLink(ctx context.Context, customerId, orderId, paymentLink string) error
}
