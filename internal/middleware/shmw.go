package middleware

import "net/http"

func SecurityHeader(next http.Handler) http.Handler {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-DNS-prefetch-Control", "off")
			w.Header().Set("X-Frame-Options", "off")
			w.Header().Set("X-Content-Type-Options", "nonsniff")

			next.ServeHTTP(w, r)
		})
}
