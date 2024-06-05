package stripe

import (
	"fmt"
	"log"

	common "github.com/naufalihsan/msvc-common"
	pb "github.com/naufalihsan/msvc-common/api"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

var (
	gatewayHttpAddress = common.EnvString("GATEWAY_HTTP_ADDR", "http://localhost:8000")
)

type Stripe struct {
}

func NewProcessor() *Stripe {
	return &Stripe{}
}

func (s *Stripe) CreatePaymentLink(order *pb.Order) (string, error) {
	log.Printf("create payment link for order %v", order)

	orderProducts := []*stripe.CheckoutSessionLineItemParams{}
	for _, orderProduct := range order.Products {
		orderProducts = append(orderProducts, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(orderProduct.PriceId),
			Quantity: stripe.Int64(int64(orderProduct.Quantity)),
		})
	}

	cancelURL := fmt.Sprintf("%s/cancel.html", gatewayHttpAddress)
	successURL := fmt.Sprintf("%s/success.html?customerId=%s&orderId=%s", gatewayHttpAddress, order.CustomerId, order.Id)

	params := &stripe.CheckoutSessionParams{
		Metadata: map[string]string{
			"orderId":    order.Id,
			"customerId": order.CustomerId,
		},
		LineItems:  orderProducts,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
	}

	result, err := session.New(params)
	if err != nil {
		return "", err
	}

	return result.URL, nil
}
