package handler

import (
	"encoding/json"
	"github.com/iscritic/hot-coffee/internal/service"
	"github.com/iscritic/hot-coffee/models"
	"log/slog"
	"net/http"
)

type OrderHandler struct {
	Service service.OrderService
}

func NewOrderHandler(mux *http.ServeMux, service service.OrderService) {
	handler := &OrderHandler{Service: service}

	mux.HandleFunc("POST /orders", handler.CreateOrder)
	mux.HandleFunc("GET /orders", handler.GetOrders)
	mux.HandleFunc("GET /orders/{id}", handler.GetOrderByID)
	mux.HandleFunc("PUT /orders/{id}", handler.UpdateOrder)
	mux.HandleFunc("DELETE /orders/{id}", handler.DeleteOrder)
	mux.HandleFunc("POST /orders/{id}/close", handler.CloseOrder)

}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	slog.Info("CREATE ORDERS")

	var order models.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
		return
	}

	err = h.Service.CreateOrder(&order)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	slog.Info("GET ORDERS")
	orders, err := h.Service.GetOrders()
	if err != nil {
		http.Error(w, `{"error": "Failed to retrieve orders"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	slog.Info("GET ORDER BY ID")

	id := r.PathValue("id")

	order, err := h.Service.GetOrderByID(id)
	if err != nil {
		http.Error(w, `{"error": "Order not found"}`, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	slog.Info("UPDATE ORDERS")

	id := r.PathValue("id")

	var order models.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
		return
	}
	order.ID = id

	err = h.Service.UpdateOrder(&order)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {

	slog.Info("DELETE ORDERS")

	id := r.PathValue("id")

	err := h.Service.DeleteOrder(id)
	if err != nil {
		http.Error(w, `{"error": "Failed to delete order"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order deleted successfully"})
}

func (h *OrderHandler) CloseOrder(w http.ResponseWriter, r *http.Request) {

	slog.Info("CLOSE ORDER")

	//TODO: id
	var id string
	err := h.Service.CloseOrder(id)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order closed successfully"})
}
