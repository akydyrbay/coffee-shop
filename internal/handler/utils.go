package handler

import (
	"encoding/json"
	"fmt"
	"frappuccino/models"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func SetBodyToJson(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	js, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		return err
	}
	w.Write(js)
	return nil
}

type ErrorResponse struct {
	Message string `json:"Error"`
}

func RespondWithJson(w http.ResponseWriter, errorResponse ErrorResponse, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}

func GetJSONRequest(w http.ResponseWriter, r *http.Request, v interface{}) {
	if r.Body == nil {
		SendResponse(w, nil, models.NotFound)
		slog.Info("Failed to read body")
	}
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		SendResponse(w, nil, models.NotFound)
		slog.Error("Failed to decode", err.Error(), "no new item to post")
	}
}

func SendResponse(w http.ResponseWriter, data interface{}, status models.Status) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status.Code)

	slog.Info("Sending response", "status", status.Code)

	if status.Code == http.StatusNoContent {
		return
	}

	if status.ErrorMessage != nil {
		errResponse := map[string]string{"error": status.ErrorMessage.Error()}
		if err := json.NewEncoder(w).Encode(errResponse); err != nil {
			slog.Error("Failed to encode error response", "error", err.Error())
		}
		slog.Error("Response error", "error", status.ErrorMessage.Error())
		return
	}

	if data == nil {
		statusResponse := map[string]interface{}{
			"code": status.Code,
		}
		if status.ErrorMessage != nil {
			statusResponse["error"] = status.ErrorMessage.Error()
		}

		if err := json.NewEncoder(w).Encode(statusResponse); err != nil {
			slog.Error("Failed to encode status response", "error", err.Error())
		}
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to encode data response", "error", err.Error())
	}
}

func ParseIndex(r *http.Request, pathLength int) (error, int) {
	pathParam := strings.Split(r.URL.Path, "/")
	if len(pathParam) != pathLength && pathLength < 2 {
		slog.Error("Failed", "wrong params", "no order posted")
		return fmt.Errorf("Invalid path"), 0
	}

	id := pathParam[2]
	ID, err := strconv.Atoi(id)
	if err != nil {
		slog.Error("Failed", err.Error(), "no order posted")
		return fmt.Errorf("Invalid id"), 0
	}

	return nil, ID
}
