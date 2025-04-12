package handler

import (
	"frappuccino/internal/service"
	"frappuccino/models"
	"log/slog"
	"net/http"
	"strconv"
)

type CustomerHandler interface {
	PostItem(w http.ResponseWriter, r *http.Request)
	GetAllItem(w http.ResponseWriter, r *http.Request)
	GetItemById(w http.ResponseWriter, r *http.Request)
	DeleteItem(w http.ResponseWriter, r *http.Request)
}

type customerHandler struct {
	service service.CustomerService
}

func NewCustomerHandler(service service.CustomerService) *customerHandler {
	return &customerHandler{service: service}
}

func (h *customerHandler) PostItem(w http.ResponseWriter, r *http.Request) {
	var item models.Customer
	GetJSONRequest(w, r, &item)
	status := h.service.CreateCustomer(item)
	SendResponse(w, nil, status)
	if status == models.Success {
		slog.Info("Customer created succesfully")
	}
}

func (h *customerHandler) GetAllItem(w http.ResponseWriter, r *http.Request) {
	code, status := h.service.GetAllCustomers(w)

	SendResponse(w, code, status)
	if status == models.Success {
		slog.Info("Customers retrieved succesfully")
	}
}

func (h *customerHandler) GetItemById(w http.ResponseWriter, r *http.Request) {
	ID := r.URL.Path[len("/customer/"):]
	id, err := strconv.Atoi(ID)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		slog.Error("Failed to read ID", err.Error(), "no new item to post")
		return
	}
	code, status := h.service.GetCustomer(w, id)

	SendResponse(w, code, status)
	if status == models.Success {
		slog.Info("Customer retrieved succesfully")
	}
}

func (h *customerHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/customer/"):]
	idForDeletion, err := strconv.Atoi(id)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		slog.Error("Failed to read ID", err.Error(), "no new item to post")
		return
	}
	status := h.service.DeleteCustomer(idForDeletion)

	SendResponse(w, nil, status)
	if status == models.Success {
		slog.Info("Inventory delete", "inventoryID", idForDeletion)
	}
}
