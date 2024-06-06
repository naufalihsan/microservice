package main

import (
	"fmt"
	"net/http"

	"github.com/naufalihsan/msvc-api-gateway/gateway"
	common "github.com/naufalihsan/msvc-common"
	pb "github.com/naufalihsan/msvc-common/api"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HttpHandler struct {
	orderGateaway gateway.OrderGateaway
}

func NewHttpHandler(orderGateaway gateway.OrderGateaway) *HttpHandler {
	return &HttpHandler{orderGateaway}
}

func (h *HttpHandler) registerRoutes(mux *http.ServeMux) {
	// handle static folder
	mux.Handle("/", http.FileServer(http.Dir("public")))

	mux.HandleFunc("POST /api/customers/{customerId}/orders", h.handleCreateOrder)
	mux.HandleFunc("GET /api/customers/{customerId}/orders/{orderId}", h.handleGetOrder)
}

func (h *HttpHandler) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	customerId := r.PathValue("customerId")
	orderId := r.PathValue("orderId")

	order, err := h.orderGateaway.GetOrder(r.Context(), customerId, orderId)

	if errStatus := status.Convert(err); errStatus != nil {
		if errStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, errStatus.Message())
			return
		}

		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.WriteJSON(w, http.StatusOK, order)
}

func (h *HttpHandler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerId := r.PathValue("customerId")

	var orderProducts []*pb.OrderProduct
	if err := common.ReadJSON(r, &orderProducts); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	tracer := otel.Tracer("http")
	ctx, span := tracer.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.RequestURI))
	defer span.End()

	order, err := h.orderGateaway.CreateOrder(ctx, &pb.CreateOrderRequest{
		CustomerId:    customerId,
		OrderProducts: orderProducts,
	})

	if errStatus := status.Convert(err); errStatus != nil {
		if errStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, errStatus.Message())
			return
		}

		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.WriteJSON(w, http.StatusOK, order)
}
