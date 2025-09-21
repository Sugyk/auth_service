package middlewares

import "net/http"

func Method(allowedMethod string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == allowedMethod {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})
}
func Post(next http.Handler) http.Handler {
	return Method(http.MethodPost, next)
}

func Get(next http.Handler) http.Handler {
	return Method(http.MethodGet, next)
}
