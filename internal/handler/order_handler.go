package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"coffee/internal/service"
	"coffee/models"
)

type OrderHandler interface {
	PostOrder(w http.ResponseWriter, r *http.Request)
	PutOrderByID(w http.ResponseWriter, r *http.Request)
	DeleteOrderByID(w http.ResponseWriter, r *http.Request)
	GetOrderByID(w http.ResponseWriter, r *http.Request)
	GetAllOrders(w http.ResponseWriter, r *http.Request)
	PostCloseOrder(w http.ResponseWriter, r *http.Request)
}

type orderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *orderHandler {
	return &orderHandler{orderService: orderService}
}

func (h *orderHandler) PostOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder models.Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed to decode", err.Error(), "no order posted")
		return
	}
	err = h.orderService.PostNewOrder(newOrder)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		slog.Error("Failed to decode", err.Error(), "no order posted")
		return
	}
	// slog.Info("order posted", "orderID")
	w.WriteHeader(http.StatusCreated)
}

func (h *orderHandler) PutOrderByID(w http.ResponseWriter, r *http.Request) {
	pathParam := strings.Split(r.URL.Path, "/")
	if len(pathParam) != 3 {
		RespondWithJson(w, ErrorResponse{Message: "Invalid path"}, http.StatusBadRequest)
		slog.Error("Failed", "wrong params", "no order posted")
		return
	}
	var newOrder models.Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
	id := pathParam[2]
	orderItem, err := h.orderService.GetOrderItemById(id)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
	if orderItem.Status == "open" {
		err = h.orderService.UpdateOrder(id, newOrder)
		if err != nil {
			if err.Error() == "not found" {
				RespondWithJson(w, ErrorResponse{Message: "Order not found"}, http.StatusNotFound)
				slog.Error("Failed", err.Error(), "no order posted")
				return
			}
			if err.Error() == "bad request" {
				RespondWithJson(w, ErrorResponse{Message: "Bad request"}, http.StatusBadRequest)
				slog.Error("Failed", err.Error(), "no order posted")
				return
			}
			RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
			slog.Error("Failed to update", err.Error(), "no order posted")
			return
		}
	} else {
		RespondWithJson(w, ErrorResponse{Message: "order closed"}, http.StatusNotFound)
		slog.Error("Failed", "no order", "order closed")
		return
	}
	slog.Info("order posted", "orderID", id)
}

func (h *orderHandler) DeleteOrderByID(w http.ResponseWriter, r *http.Request) {
	pathParam := strings.Split(r.URL.Path, "/")
	if len(pathParam) != 3 {
		RespondWithJson(w, ErrorResponse{Message: "Invalid path"}, http.StatusBadRequest)
		slog.Error("Failed", "wrong params", "no order posted")
		return
	}
	id := pathParam[2]
	err := h.orderService.DeleteOrder(id)
	if err != nil {
		if err.Error() == "not found" {
			RespondWithJson(w, ErrorResponse{Message: "Order not found"}, http.StatusNotFound)
			slog.Error("Failed", err.Error(), "no order posted")
			return
		}
		if err.Error() == "bad request" {
			RespondWithJson(w, ErrorResponse{Message: "Bad request"}, http.StatusBadRequest)
			slog.Error("Failed", err.Error(), "no order posted")
			return
		}
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed", err.Error(), "no order posted")
	}
	slog.Info("order posted", "orderID", id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *orderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/orders/"):]
	orderItem, err := h.orderService.GetOrderItemById(id)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
	err = setBodyToJson(w, orderItem)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
	slog.Info("order posted", "orderID", id)
}

func (h *orderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orderItems, err := h.orderService.GetOrderItem()
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}

	jsonData, err := json.MarshalIndent(orderItems, "", "  ")
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
	slog.Info("order posted", "orderID", orderItems)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *orderHandler) PostCloseOrder(w http.ResponseWriter, r *http.Request) {
	pathParam := strings.Split(r.URL.Path, "/")
	if len(pathParam) != 4 {
		RespondWithJson(w, ErrorResponse{Message: "Invalid path"}, http.StatusBadRequest)
		slog.Error("Failed", "wrong params", "no order posted")
		return
	}
	id := pathParam[2]
	if err := h.orderService.UpdateOrderStatus(id); err != nil {
		if err.Error() == "order is already closed" {
			RespondWithJson(w, ErrorResponse{Message: "Order is already closed"}, http.StatusNotFound)
			slog.Error("Failed", err.Error(), "order is already closed")
			return
		}
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
	slog.Info("order posted", "orderID", id)
}
