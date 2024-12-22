package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"coffee/internal/service"
	"coffee/models"
)

type MenuHandler interface {
	PostMenuHandler(w http.ResponseWriter, r *http.Request)
	GetAllMenuHandler(w http.ResponseWriter, r *http.Request)
	GetMenuItemHandler(w http.ResponseWriter, r *http.Request)
	PutMenuHandler(w http.ResponseWriter, r *http.Request)
	DeleteMenuHandler(w http.ResponseWriter, r *http.Request)
}

type menuHandler struct {
	menuService service.MenuServiceInterface
}

func NewMenuHandler(menuService service.MenuServiceInterface) *menuHandler {
	return &menuHandler{menuService: menuService}
}

func (h *menuHandler) PostMenuHandler(w http.ResponseWriter, r *http.Request) {
	var newMenuitem models.MenuItem
	json.NewDecoder(r.Body).Decode(&newMenuitem)
	err := h.menuService.AddMenuItem(newMenuitem)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		slog.Error("Failed to AddMenuItem", err.Error(), "no menu posted")
		return
	}
	slog.Info("menu posted", "menuID", newMenuitem.ID)
	w.WriteHeader(http.StatusCreated)
}

func (h *menuHandler) GetAllMenuHandler(w http.ResponseWriter, r *http.Request) {
	menuItems, err := h.menuService.GetAllMenuItems()
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed to GetAllMenuItems", err.Error(), "no menu posted")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonMenus, err := json.MarshalIndent(menuItems, "", "  ")
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed to MarshalIndent", err.Error(), "no menu posted")
		return
	}
	slog.Info("menu posted", "menuID", menuItems)
	w.Write(jsonMenus)
}

func (h *menuHandler) GetMenuItemHandler(w http.ResponseWriter, r *http.Request) {
	pathParam := strings.Split(r.URL.Path, "/")
	if len(pathParam) != 3 {
		RespondWithJson(w, ErrorResponse{Message: "Invalid path"}, http.StatusBadRequest)
		slog.Error("Failed to get input", "wrong input", "no menu posted")
		return
	}
	id := pathParam[2]
	menuItem, err := h.menuService.GetMenuItemById(id)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed to GetMenuItemById", err.Error(), "no menu posted")
		return
	}
	err = setBodyToJson(w, menuItem)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed to setBodyToJson", err.Error(), "no menu posted")
		return
	}
	slog.Info("menu posted", "menuID", menuItem.ID)
}

func (h *menuHandler) PutMenuHandler(w http.ResponseWriter, r *http.Request) {
	var menuItem models.MenuItem
	id := r.URL.Path[len("/menu/"):]
	err := json.NewDecoder(r.Body).Decode(&menuItem)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		slog.Error("Failed to decode", err.Error(), "no menu posted")
		return
	}
	if menuItem.ID != id {
		RespondWithJson(w, ErrorResponse{Message: "Menu ID conflict"}, http.StatusBadRequest)
	}
	err = h.menuService.UpdateMenuItem(menuItem)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed to UpdateMenuItem", err.Error(), "no menu posted")
		return
	}
	slog.Info("menu posted", "menuID", menuItem.ID)
}

func (h *menuHandler) DeleteMenuHandler(w http.ResponseWriter, r *http.Request) {
	pathParam := strings.Split(r.URL.Path, "/")
	if len(pathParam) != 3 {
		RespondWithJson(w, ErrorResponse{Message: "Invalid path"}, http.StatusBadRequest)
		slog.Error("Failed to get", "wrong input", "no menu posted")
		return
	}
	id := pathParam[2]
	err := h.menuService.DeleteMenuItemById(id)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed to DeleteMenuItemById", err.Error(), "no menu posted")
		return
	}
	slog.Info("menu posted", "menuID", id)
	w.WriteHeader(http.StatusNoContent)
}
