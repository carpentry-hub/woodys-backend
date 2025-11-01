package middlewares

import "net/http"

func JsonContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		// Establecer la cabecera para cada respuesta
		w.Header().Set("Content-Type", "application/json")
		
		// Llama al siguiente handler en la cadena
		next.ServeHTTP(w, r)
	})
}