package handler

import (
	"encoding/json"
	"frappuccino/internal/service"
	"frappuccino/models"
	"log/slog"
	"net/http"
	"strconv"
)

type AggragationHandler interface {
	GetAllSales(w http.ResponseWriter, r *http.Request)
	GetPopularSales(w http.ResponseWriter, r *http.Request)
	GetOrderByDate(w http.ResponseWriter, r *http.Request)
	FullTextSearchReport(w http.ResponseWriter, r *http.Request)
	OrderedItemsByPeriod(w http.ResponseWriter, r *http.Request)
	GetLeftOvers(w http.ResponseWriter, r *http.Request)
}

type aggragationHandler struct {
	aggragationService service.AggragationService
}

func NewAggragationHandler(aggragation service.AggragationService) *aggragationHandler {
	return &aggragationHandler{aggragationService: aggragation}
}

func (h *aggragationHandler) GetAllSales(w http.ResponseWriter, r *http.Request) {
	salesAmount, _ := h.aggragationService.GetTotalSales()
	var total models.TotalSales
	total.Sales = salesAmount
	err := SetBodyToJson(w, total)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed", err.Error(), "no total sales to post")
	}
}

func (h *aggragationHandler) GetPopularSales(w http.ResponseWriter, r *http.Request) {
	list, err := h.aggragationService.GetPopularMenuItems()
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed", err.Error(), "no popular items to post")
	}
	jsonData, err := json.MarshalIndent(list, "", "   ")
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed", err.Error(), "no popular items to post")
	}
	slog.Info("popular sales posted", "popular", list)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (h *aggragationHandler) GetOrderByDate(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	orderItems, err := h.aggragationService.GetOrderQuantityByDate(startDate, endDate)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
	result := make(map[string]int)
	for _, item := range orderItems {
		result[item.Name] = int(item.Quantity)
	}
	err = SetBodyToJson(w, result)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
	// slog.Info("order posted", "orderID", id)
}

func (h *aggragationHandler) FullTextSearchReport(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	filter := r.URL.Query().Get("filter")
	minPrice := r.URL.Query().Get("minPrice")
	maxPrice := r.URL.Query().Get("maxPrice")

	orderItems, err := h.aggragationService.FullSearchReport(q, filter, minPrice, maxPrice)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
	// result := make(map[string]int)
	// for _, item := range orderItems {
	// 	result[item.Name] = int(item.Quantity)
	// }
	err = SetBodyToJson(w, orderItems)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
	slog.Info("success")
}

func (h *aggragationHandler) OrderedItemsByPeriod(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	var month, year string
	switch period {
	case "day":
		month = r.URL.Query().Get("month")
		if month == "" {
			month = "january"
		}
		m, err := h.aggragationService.GetOrderedItemsByDay(month)
		if err != nil {
			RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
			slog.Error("Failed", err.Error(), "no order posted")
			return
		}
		var items []map[string]int
		for _, item := range m.OrderedItems {
			dayMap := map[string]int{
				item.Day: item.Quantity,
			}
			items = append(items, dayMap)
		}

		response := models.OrderByDayRequest2{
			Period:       m.Period,
			Month:        m.Month,
			OrderedItems: items,
		}
		if err := SetBodyToJson(w, response); err != nil {
			RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
			slog.Error("Failed to send response", "error", err.Error())
			return
		}
	case "month":
		year = r.URL.Query().Get("year")
		if year == "" {
			year = "2024"
		}
		y, err := h.aggragationService.GetOrderedItemsByMonth(year)
		if err != nil {
			RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
			slog.Error("Failed", err.Error(), "no order posted")
			return
		}
		var items []map[string]int
		for _, item := range y.OrderedItems {
			dayMap := map[string]int{
				item.Month: item.Quantity,
			}
			items = append(items, dayMap)
		}

		response := models.OrderByDayRequest2{
			Period:       y.Period,
			Month:        y.Year,
			OrderedItems: items,
		}
		if err := SetBodyToJson(w, response); err != nil {
			RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
			slog.Error("Failed to send response", "error", err.Error())
			return
		}
	default:
		RespondWithJson(w, ErrorResponse{Message: "wrong input"}, http.StatusInternalServerError)
		slog.Info("Failed")
		return
	}
	// slog.Info("order posted", "orderID", id)
}

func (h *aggragationHandler) GetLeftOvers(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sortBy")
	pageStr := r.URL.Query().Get("page")
	var page int
	if pageStr == "" {
		page = 1
	} else {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
			slog.Error("Failed", err.Error(), "no order posted")
			return
		}
	}
	pageSizeStr := r.URL.Query().Get("pageSize")
	var pageSize int
	if pageSizeStr == "" {
		pageSize = 10
	} else {
		var err error
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil || pageSize <= 0 {
			RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
			slog.Error("Failed", err.Error(), "no order posted")
			return
		}
	}
	if sortBy == "" {
		sortBy = "price"
	} else if sortBy != "price" && sortBy != "quantity" {
		RespondWithJson(w, ErrorResponse{Message: ""}, http.StatusNotFound)
		// slog.Error("Failed", "no order posted")
		return
	}
	response, err := h.aggragationService.GetLeftOvers(sortBy, page, pageSize)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
	if err := SetBodyToJson(w, response); err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed to send response", "error", err.Error())
		return
	}
	slog.Info("Leftovers retrieved successfully")
}
