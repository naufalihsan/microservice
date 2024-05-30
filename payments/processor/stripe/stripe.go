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

	successURL := fmt.Sprintf("%s/payments/success.html", gatewayHttpAddress)

	params := &stripe.CheckoutSessionParams{
		LineItems:  orderProducts,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(successURL),
	}

	result, err := session.New(params)
	if err != nil {
		return "", err
	}

	return result.URL, nil
}
