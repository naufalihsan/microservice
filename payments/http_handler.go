package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	common "github.com/naufalihsan/msvc-common"
	pb "github.com/naufalihsan/msvc-common/api"
	"github.com/naufalihsan/msvc-common/broker"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
)

type HttpHandler struct {
	channel *amqp.Channel
}

func NewHttpHandler(channel *amqp.Channel) *HttpHandler {
	return &HttpHandler{channel}
}

func (h *HttpHandler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/webhook", h.handleStripeWebhook)
}

func (h *HttpHandler) handleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	fmt.Fprintf(os.Stdout, "Got body: %s\n", body)

	// Pass the request body and Stripe-Signature header to ConstructEvent, along with the webhook signing key
	// You can find your endpoint's secret in your webhook settings
	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), stripeEndpointSecret)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	// Handle the checkout.session.completed event
	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		order := &pb.Order{
			Id:          session.Metadata["orderId"],
			CustomerId:  session.Metadata["customerId"],
			Status:      common.OrderStatusPaid,
			PaymentLink: session.PaymentLink.ID,
		}

		jsonOrder, err := json.Marshal(order)
		if err != nil {
			log.Fatal(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if session.PaymentStatus == "paid" {
			h.channel.PublishWithContext(ctx, broker.OrderPaidEvent, "", false, false, amqp.Publishing{
				ContentType: "application/json",
				Body:        jsonOrder,
				Priority:    amqp.Persistent,
			})
		}
	}

	w.WriteHeader(http.StatusOK)
}
