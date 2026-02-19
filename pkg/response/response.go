package response

import (
	"encoding/json"
	"net/http"
)

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

type PaginatedResponse struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type ErrorResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func WriteSuccess(w http.ResponseWriter, data interface{}, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})

	if err != nil {
		WriteError(w, "Failed to encode success response", http.StatusInternalServerError)
		return
	}
}

func WriteError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(ErrorResponse{
		Success: false,
		Message: message,
	})

	if err != nil {
		WriteError(w, "Failed to encode error response", http.StatusInternalServerError)
		return
	}
}

func WriteValidationErrors(w http.ResponseWriter, errors map[string]string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(ErrorResponse{
		Success: false,
		Message: "Validation failed",
		Errors:  errors,
	})

	if err != nil {
		WriteError(w, "Failed to encode validation error response", http.StatusInternalServerError)
		return
	}
}

func WritePaginated(w http.ResponseWriter, data interface{}, meta PaginationMeta, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(PaginatedResponse{
		Data: data,
		Meta: meta,
	})

	if err != nil {
		WriteError(w, "Failed to encode paginated response", http.StatusInternalServerError)
		return
	}
}
