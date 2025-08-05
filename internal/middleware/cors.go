package middleware

import (
	"net/http"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig returns default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
			http.MethodPatch,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-Request-ID",
		},
		ExposeHeaders: []string{
			"X-Request-ID",
		},
		AllowCredentials: false,
		MaxAge:           86400, // 24 hours
	}
}

// CORS returns a middleware that handles CORS
func CORS(config ...CORSConfig) func(http.Handler) http.Handler {
	cfg := DefaultCORSConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Set Access-Control-Allow-Origin
			if len(cfg.AllowOrigins) == 1 && cfg.AllowOrigins[0] == "*" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else if origin != "" && isOriginAllowed(origin, cfg.AllowOrigins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			// Set Access-Control-Allow-Methods
			if len(cfg.AllowMethods) > 0 {
				w.Header().Set("Access-Control-Allow-Methods", joinStrings(cfg.AllowMethods, ", "))
			}

			// Set Access-Control-Allow-Headers
			if len(cfg.AllowHeaders) > 0 {
				w.Header().Set("Access-Control-Allow-Headers", joinStrings(cfg.AllowHeaders, ", "))
			}

			// Set Access-Control-Expose-Headers
			if len(cfg.ExposeHeaders) > 0 {
				w.Header().Set("Access-Control-Expose-Headers", joinStrings(cfg.ExposeHeaders, ", "))
			}

			// Set Access-Control-Allow-Credentials
			if cfg.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Set Access-Control-Max-Age
			if cfg.MaxAge > 0 {
				w.Header().Set("Access-Control-Max-Age", intToString(cfg.MaxAge))
			}

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// isOriginAllowed checks if the origin is in the allowed origins list
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == "*" || allowedOrigin == origin {
			return true
		}
	}
	return false
}

// joinStrings joins a slice of strings with a separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// intToString converts an integer to string
func intToString(i int) string {
	return string(rune(i + '0'))
}
