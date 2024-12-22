package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"coffee/internal/service"
	"coffee/models"
)

type AggragationHandler interface {
	GetAllSales(w http.ResponseWriter, r *http.Request)
	GetPopularSales(w http.ResponseWriter, r *http.Request)
}

type aggragationHandler struct {
	aggragationService service.AggragationService
}

func NewAggragationHandler(service service.AggragationService) *aggragationHandler {
	return &aggragationHandler{aggragationService: service}
}

func (h *aggragationHandler) GetAllSales(w http.ResponseWriter, r *http.Request) {
	salesAmount, err := h.aggragationService.GetTotalSales()
	if err != nil {
		message := err.Error()
		if len(message) > 20 {
			if message[:20] == "menu item not found:" {
				RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
				slog.Error("Failed", err.Error(), "no total sales to post")
			}
		}
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed", err.Error(), "no total sales to post")
	}
	jsonData, err := json.MarshalIndent(models.TotalSales{Sales: salesAmount}, "", "   ")
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed", err.Error(), "no total sales to post")
	}
	slog.Info("total sales posted", "total", salesAmount)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (h *aggragationHandler) GetPopularSales(w http.ResponseWriter, r *http.Request) {
	list, err := h.aggragationService.GetPopularMenuItems()
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed", err.Error(), "no total sales to post")
	}
	jsonData, err := json.MarshalIndent(list, "", "   ")
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed", err.Error(), "no total sales to post")
	}
	slog.Info("popular sales posted", "popular", list)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
