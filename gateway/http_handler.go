package main

import (
	"log"
	"net/http"
)

type HttpHandler struct {
}

func NewHttpHandler() *HttpHandler {
	return &HttpHandler{}
}

func (h *HttpHandler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.handleCreateOrder)
}

func (h *HttpHandler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	log.Print("air hot reload")
}
