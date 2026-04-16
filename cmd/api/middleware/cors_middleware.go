package middleware

import "net/http"

func CorsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		println("token", token)
		origin := r.Header.Get("Origin")

		allowed := map[string]bool{
			"http://127.0.0.1":      true,
			"http://localhost":      true,
			"http://120.0.1.1:80":   true,
			"http://localhost:80":   true,
			"http://localhost:4200": true,
			"http://192.168.0.24":   true,
		}

		if allowed[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
