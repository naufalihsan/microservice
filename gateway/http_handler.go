package main

import (
	"net/http"

	common "github.com/naufalihsan/msvc-common"
	pb "github.com/naufalihsan/msvc-common/api"
)

type HttpHandler struct {
	client pb.OrderServiceClient
}

func NewHttpHandler(client pb.OrderServiceClient) *HttpHandler {
	return &HttpHandler{client}
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

	h.client.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerId:    customerId,
		OrderProducts: orderProducts,
	})
}
