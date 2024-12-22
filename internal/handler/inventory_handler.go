package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"coffee/internal/service"
	"coffee/models"
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
	var newInventoryItem models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&newInventoryItem); err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed to decode", err.Error(), "no new item to post")
		return
	}
	if err := h.inventoryService.AddInventoryItem(newInventoryItem); err != nil {
		if err.Error() == "item already exists" {
			RespondWithJson(w, ErrorResponse{Message: "item already exists"}, http.StatusConflict)
			slog.Error("Item already exists")
			return
		}

		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed AddInventoryItem", err.Error(), "no new item to post")
		return
	}
	slog.Info("Inventory posted", "inventoryID", newInventoryItem.IngredientID)
	w.WriteHeader(http.StatusCreated)
}

func (h *inventoryHandler) GetAllItem(w http.ResponseWriter, r *http.Request) {
	inventoryItems, err := h.inventoryService.GetInventoryItem()
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed to get", err.Error(), "no new item to post")
		return
	}

	jsonData, err := json.MarshalIndent(inventoryItems, "", "  ")
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed to MarshalIndent", err.Error(), "no new item to post")
		return
	}
	slog.Info("Inventory got", "inventoryID", inventoryItems)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *inventoryHandler) GetItemById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/inventory/"):]

	inventoryItem, err := h.inventoryService.GetInventoryItemById(id)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed to get", err.Error(), "no new item to post")
		return
	}
	err = setBodyToJson(w, inventoryItem)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed to get", err.Error(), "no new item to post")
		return
	}
	slog.Info("Inventory got", "inventoryID", inventoryItem.IngredientID)
}

func (h *inventoryHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	idForDeletion := r.URL.Path[len("/inventory/"):]
	if err := h.inventoryService.DeleteInventoryItem(idForDeletion); err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed to MarshalIndent", err.Error(), "no new item to post")
		return
	}
	slog.Info("Inventory delete", "inventoryID", idForDeletion)
	w.WriteHeader(http.StatusNoContent)
}

func (h *inventoryHandler) PutItem(w http.ResponseWriter, r *http.Request) {
	var inventoryItem models.InventoryItem
	err := json.NewDecoder(r.Body).Decode(&inventoryItem)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		slog.Error("Failed to decode", err.Error(), "no new item to post")
		return
	}
	err = h.inventoryService.UpdateInventoryItem(inventoryItem)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed to MarshalIndent", err.Error(), "no new item to post")
		return
	}
	slog.Info("Inventory put", "inventoryID", inventoryItem.IngredientID)
}
