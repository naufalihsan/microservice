package main

import (
	"net/http"

	"github.com/naufalihsan/msvc-api-gateway/gateway"
	common "github.com/naufalihsan/msvc-common"
	pb "github.com/naufalihsan/msvc-common/api"
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
	mux.HandleFunc("POST /api/customers/{customerId}/orders", h.handleCreateOrder)
}

func (h *HttpHandler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerId := r.PathValue("customerId")

	var orderProducts []*pb.OrderProduct
	if err := common.ReadJSON(r, &orderProducts); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.orderGateaway.CreateOrder(r.Context(), &pb.CreateOrderRequest{
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
