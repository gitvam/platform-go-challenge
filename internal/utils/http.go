package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// swagger:model
type SuccessResponse struct {
	Status string      `json:"status"` 
	Data   interface{} `json:"data"`   
}

// swagger:model
type ErrorResponse struct {
	Status  string `json:"status"`  
	Message string `json:"message"` 
}

func WriteJSONError(w http.ResponseWriter, msg string, status int) {
	log.Printf("[ERROR] %d - %s", status, msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Status:  "error",
		Message: msg,
	})
}
