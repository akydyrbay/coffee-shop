package handler

import (
	"fmt"
	"frappuccino/internal/service"
	"frappuccino/models"
	"log/slog"
	"net/http"
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
	var newMenuitem []models.MenuItem
	GetJSONRequest(w, r, &newMenuitem)
	status := h.menuService.Create(newMenuitem)
	SendResponse(w, nil, status)
	slog.Info("menu posted")
}

func (h *menuHandler) GetAllMenuHandler(w http.ResponseWriter, r *http.Request) {
	menuItems, status := h.menuService.ReadAll()
	SendResponse(w, menuItems, status)
	slog.Info("menu got")
}

func (h *menuHandler) GetMenuItemHandler(w http.ResponseWriter, r *http.Request) {
	err, ID := ParseIndex(r, 3)
	if err != nil {
		SendResponse(w, nil, models.Status{ErrorMessage: fmt.Errorf("failed to parse index: %w", err), Code: 400})
		return
	}
	menuItem, status := h.menuService.Read(ID)
	SendResponse(w, menuItem, status)
	slog.Info("menu posted", "menuID", ID)
}

func (h *menuHandler) PutMenuHandler(w http.ResponseWriter, r *http.Request) {
	err, ID := ParseIndex(r, 3)
	if err != nil {
		SendResponse(w, nil, models.Status{ErrorMessage: fmt.Errorf("failed to read ID: %w", err), Code: http.StatusBadRequest})
		slog.Info("Failed to read ID")
		return
	}
	var menuItem models.MenuItem
	GetJSONRequest(w, r, &menuItem)
	menuItem.ID = ID
	status := h.menuService.Update(menuItem)
	SendResponse(w, nil, status)
	slog.Info("menu posted", "menuID", ID)
}

func (h *menuHandler) DeleteMenuHandler(w http.ResponseWriter, r *http.Request) {
	err, ID := ParseIndex(r, 3)

	if err != nil {
		SendResponse(w, nil, models.Status{ErrorMessage: fmt.Errorf("failed to read ID: %w", err), Code: http.StatusBadRequest})
		slog.Info("Failed to read ID")
		return
	}
	status := h.menuService.Delete(ID)
	SendResponse(w, nil, status)
	slog.Info("menu posted", "menuID", ID)
}
