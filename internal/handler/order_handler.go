package handler

import (
	"fmt"
	"frappuccino/internal/service"
	"frappuccino/models"
	"log/slog"
	"net/http"
)

type OrderHandler interface {
	PostOrder(w http.ResponseWriter, r *http.Request)
	PutOrderByID(w http.ResponseWriter, r *http.Request)
	DeleteOrderByID(w http.ResponseWriter, r *http.Request)
	GetOrderByID(w http.ResponseWriter, r *http.Request)
	GetAllOrders(w http.ResponseWriter, r *http.Request)
	// GetOrderByDate(w http.ResponseWriter, r *http.Request)
}

type orderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *orderHandler {
	return &orderHandler{orderService: orderService}
}

func (h *orderHandler) PostOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder []models.Order
	GetJSONRequest(w, r, &newOrder)
	status := h.orderService.Create(newOrder)
	SendResponse(w, nil, status)
	slog.Info("order posted")
}

func (h *orderHandler) PutOrderByID(w http.ResponseWriter, r *http.Request) {
	err, ID := ParseIndex(r, 3)
	if err != nil {
		SendResponse(w, nil, models.Status{ErrorMessage: fmt.Errorf("failed to read ID: %w", err), Code: http.StatusBadRequest})
		slog.Info("Failed to read ID")
		return
	}
	var orderItem models.Order
	GetJSONRequest(w, r, &orderItem)
	orderItem.ID = ID
	status := h.orderService.Update(orderItem)
	SendResponse(w, nil, status)
	slog.Info("order updated", "orderID", ID)
}

func (h *orderHandler) DeleteOrderByID(w http.ResponseWriter, r *http.Request) {
	err, ID := ParseIndex(r, 3)
	if err != nil {
		SendResponse(w, nil, models.Status{ErrorMessage: fmt.Errorf("failed to read ID: %w", err), Code: http.StatusBadRequest})
		slog.Info("Failed to read ID")
		return
	}

	status := h.orderService.Delete(ID)
	SendResponse(w, nil, status)
	slog.Info("order deleted", "orderID", ID)
}

func (h *orderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	err, ID := ParseIndex(r, 3)
	if err != nil {
		SendResponse(w, nil, models.Status{ErrorMessage: fmt.Errorf("failed to parse index: %w", err), Code: 400})
		return
	}
	orderItem, status := h.orderService.Read(ID)
	SendResponse(w, orderItem, status)
	slog.Info("order get by id", "orderID", ID)
}

func (h *orderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orderItems, status := h.orderService.ReadAll()
	SendResponse(w, orderItems, status)
	slog.Info("order getall")
}
