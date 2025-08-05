package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/carpentry-hub/woodys-backend/internal/middleware"
	"github.com/carpentry-hub/woodys-backend/internal/services"
)

// Handlers holds all HTTP handlers and their dependencies
type Handlers struct {
	services *services.Services
}

// NewHandlers creates a new handlers instance
func NewHandlers(services *services.Services) *Handlers {
	return &Handlers{
		services: services,
	}
}

// APIResponse represents a standard API response format
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Message   string      `json:"message,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// writeJSONResponse writes a JSON response with the standard format
func (h *Handlers) writeJSONResponse(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}, errorMsg, message string) {
	response := APIResponse{
		Success:   statusCode < 400,
		Data:      data,
		Error:     errorMsg,
		Message:   message,
		RequestID: middleware.GetRequestIDFromContext(r.Context()),
		Timestamp: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// writeErrorResponse writes an error response
func (h *Handlers) writeErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, errorMsg, message string) {
	h.writeJSONResponse(w, r, statusCode, nil, errorMsg, message)
}

// writeSuccessResponse writes a success response
func (h *Handlers) writeSuccessResponse(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}, message string) {
	h.writeJSONResponse(w, r, statusCode, data, "", message)
}

// parseIDFromPath extracts and validates an ID from the URL path
func (h *Handlers) parseIDFromPath(r *http.Request, paramName string) (int64, error) {
	vars := mux.Vars(r)
	idStr, exists := vars[paramName]
	if !exists {
		return 0, &ValidationError{Field: paramName, Message: "ID parameter is required"}
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return 0, &ValidationError{Field: paramName, Message: "Invalid ID format"}
	}

	return id, nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// Home handles the root endpoint
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"service": "Woody's Backend API",
		"version": "1.0.0",
		"status":  "running",
	}
	h.writeSuccessResponse(w, r, http.StatusOK, response, "API is running")
}

// HealthCheck handles health check requests
func (h *Handlers) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// TODO: Add actual health checks (database connectivity, etc.)
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"checks": map[string]string{
			"database": "ok",
			"memory":   "ok",
		},
	}
	h.writeSuccessResponse(w, r, http.StatusOK, response, "Service is healthy")
}

// getPaginationParams extracts pagination parameters from query string
func (h *Handlers) getPaginationParams(r *http.Request) (limit, offset int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit = 50 // default limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
			if limit > 100 {
				limit = 100 // max limit
			}
		}
	}

	offset = 0 // default offset
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	return limit, offset
}

// validateRequestBody validates and decodes JSON request body
func (h *Handlers) validateRequestBody(r *http.Request, dest interface{}) error {
	if r.Body == nil {
		return &ValidationError{Field: "body", Message: "Request body is required"}
	}

	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return &ValidationError{Field: "body", Message: "Invalid JSON format: " + err.Error()}
	}

	return nil
}

// handleInternalError logs and returns internal server error
func (h *Handlers) handleInternalError(w http.ResponseWriter, r *http.Request, err error, message string) {
	// TODO: Add proper logging
	h.writeErrorResponse(w, r, http.StatusInternalServerError, "Internal Server Error", message)
}

// handleValidationError returns validation error response
func (h *Handlers) handleValidationError(w http.ResponseWriter, r *http.Request, err error) {
	if valErr, ok := err.(*ValidationError); ok {
		h.writeErrorResponse(w, r, http.StatusBadRequest, "Validation Error", valErr.Message)
	} else {
		h.writeErrorResponse(w, r, http.StatusBadRequest, "Bad Request", err.Error())
	}
}

// handleNotFoundError returns not found error response
func (h *Handlers) handleNotFoundError(w http.ResponseWriter, r *http.Request, resource string) {
	h.writeErrorResponse(w, r, http.StatusNotFound, "Not Found", resource+" not found")
}
