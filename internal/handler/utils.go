package handler

import (
	"encoding/json"
	"net/http"
)

func setBodyToJson(w http.ResponseWriter, data interface{}) error {
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
