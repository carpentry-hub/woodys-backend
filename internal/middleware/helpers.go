package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"
)

// Context keys for request data
type contextKey string

const (
	requestIDKey contextKey = "request_id"
	userIDKey    contextKey = "user_id"
)

// generateRequestID generates a unique request ID
func generateRequestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to a simple implementation if crypto/rand fails
		return "req_" + hex.EncodeToString([]byte("fallback"))
	}
	return "req_" + hex.EncodeToString(bytes)
}

// setRequestID sets the request ID in the context
func setRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// getRequestID gets the request ID from the request context
func getRequestID(r *http.Request) string {
	if requestID, ok := r.Context().Value(requestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// GetRequestIDFromContext gets the request ID from context
func GetRequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// setUserID sets the user ID in the context
func setUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserIDFromContext gets the user ID from context
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	if userID, ok := ctx.Value(userIDKey).(int64); ok {
		return userID, true
	}
	return 0, false
}

// SetUserIDInContext is a helper to set user ID in request context
func SetUserIDInContext(r *http.Request, userID int64) *http.Request {
	ctx := setUserID(r.Context(), userID)
	return r.WithContext(ctx)
}

// WriteJSONError writes a JSON error response
func WriteJSONError(w http.ResponseWriter, statusCode int, errorMsg, message string, requestID string) {
	errorResponse := ErrorResponse{
		Error:     errorMsg,
		Message:   message,
		RequestID: requestID,
		Timestamp: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}

// Chain chains multiple middleware functions
func Chain(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}
		return handler
	}
}

// ValidateContentType validates that the request has the correct content type
func ValidateContentType(contentType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
				if r.Header.Get("Content-Type") != contentType {
					WriteJSONError(w, http.StatusUnsupportedMediaType,
						"Unsupported Media Type",
						"Content-Type must be "+contentType,
						getRequestID(r))
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// MethodNotAllowed returns a 405 Method Not Allowed response
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	WriteJSONError(w, http.StatusMethodNotAllowed,
		"Method Not Allowed",
		"The HTTP method "+r.Method+" is not allowed for this endpoint",
		getRequestID(r))
}

// NotFound returns a 404 Not Found response
func NotFound(w http.ResponseWriter, r *http.Request) {
	WriteJSONError(w, http.StatusNotFound,
		"Not Found",
		"The requested resource was not found",
		getRequestID(r))
}
