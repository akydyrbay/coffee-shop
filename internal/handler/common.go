package handler

import (
	"encoding/json"
	"fmt"
	"frappuccino/internal/service"
	"frappuccino/models"
	"log/slog"
	"net/http"
)

type Handler interface {
	Close(w http.ResponseWriter, r *http.Request)
	BatchProcessOrders(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service service.Service
}

func NewHandler(service service.Service) *handler {
	return &handler{service: service}
}

func (h *handler) Close(w http.ResponseWriter, r *http.Request) {
	err, ID := ParseIndex(r, 4)
	if err != nil {
		SendResponse(w, nil, models.Status{ErrorMessage: fmt.Errorf("failed to read ID: %w", err), Code: http.StatusBadRequest})
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}

	status := h.service.CloseOrdersByID(ID)
	SendResponse(w, nil, status)
}

func (h *handler) BatchProcessOrders(w http.ResponseWriter, r *http.Request) {
	var request models.BatchProcessRequest

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		SendResponse(w, nil, models.Status{ErrorMessage: fmt.Errorf("failed to read ID: %w", err), Code: http.StatusBadRequest})
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}

	// Call service layer
	response, err := h.service.ProcessBatchOrders(request)
	if err != nil {
		SendResponse(w, nil, models.Status{ErrorMessage: fmt.Errorf("failed to read ID: %w", err), Code: http.StatusBadRequest})
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}

	// Send response
	if err := SetBodyToJson(w, response); err != nil {
		SendResponse(w, nil, models.Status{ErrorMessage: fmt.Errorf("failed to read ID: %w", err), Code: http.StatusBadRequest})
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}

	slog.Info("Order batch process has been successfully completed")
}
