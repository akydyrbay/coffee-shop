package handler

import (
	"fmt"
	"frappuccino/internal/service"
	"frappuccino/models"
	"log/slog"
	"net/http"
)

type InventoryHandler interface {
	PostItem(w http.ResponseWriter, r *http.Request)
	GetAllItem(w http.ResponseWriter, r *http.Request)
	GetItemById(w http.ResponseWriter, r *http.Request)
	DeleteItem(w http.ResponseWriter, r *http.Request)
	PutItem(w http.ResponseWriter, r *http.Request)
}

type inventoryHandler struct {
	inventoryService service.InventoryService
}

func NewInventoryHandler(inventoryService service.InventoryService) *inventoryHandler {
	return &inventoryHandler{inventoryService: inventoryService}
}

func (h *inventoryHandler) PostItem(w http.ResponseWriter, r *http.Request) {
	var newInventoryItem []models.InventoryItem
	GetJSONRequest(w, r, &newInventoryItem)

	status := h.inventoryService.Create(newInventoryItem)
	SendResponse(w, nil, status)
	if status != models.BadRequest {
		slog.Info("Inventory posted")
	}
}

func (h *inventoryHandler) GetAllItem(w http.ResponseWriter, r *http.Request) {
	inventoryItems, status := h.inventoryService.ReadAll()
	SendResponse(w, inventoryItems, status)
	if status != models.BadRequest {
		slog.Info("Inventory got")
	}
}

func (h *inventoryHandler) GetItemById(w http.ResponseWriter, r *http.Request) {
	err, ID := ParseIndex(r, 3)
	if err != nil {
		SendResponse(w, nil, models.Status{ErrorMessage: fmt.Errorf("failed to read ID: %w", err), Code: http.StatusBadRequest})
		slog.Error("Failed to read ID")
		return
	}
	inventoryItem, status := h.inventoryService.Read(ID)
	SendResponse(w, inventoryItem, status)
	if status != models.NotFound {
		slog.Info("Inventory got", "inventoryID", inventoryItem.ID)
	}
}

func (h *inventoryHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	err, ID := ParseIndex(r, 3)

	if err != nil {
		SendResponse(w, nil, models.Status{ErrorMessage: fmt.Errorf("failed to read ID: %w", err), Code: http.StatusBadRequest})
		slog.Error("Failed to read ID")
		return
	}
	status := h.inventoryService.Delete(ID)
	SendResponse(w, nil, status)

	slog.Info("Inventory delete", "inventoryID", ID)
}

func (h *inventoryHandler) PutItem(w http.ResponseWriter, r *http.Request) {
	err, ID := ParseIndex(r, 3)
	if err != nil {
		SendResponse(w, nil, models.Status{ErrorMessage: fmt.Errorf("failed to read ID: %w", err), Code: http.StatusBadRequest})
		slog.Error("Failed to read ID")
		return
	}
	var inventoryItem models.InventoryItem
	GetJSONRequest(w, r, &inventoryItem)
	inventoryItem.ID = ID
	status := h.inventoryService.Update(inventoryItem)
	SendResponse(w, nil, status)
	slog.Info("Inventory put", "inventoryID", inventoryItem.ID)
}
