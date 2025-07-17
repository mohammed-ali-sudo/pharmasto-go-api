package middleware

import "net/http"

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// إعداد رؤوس CORS
		w.Header().Set("Access-Control-Allow-Origin", "*") // أو ضع موقعك بدل *
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// إذا كان الطلب من نوع OPTIONS (preflight)، نرد وننهي بدون تمرير للـ handler
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// استدعاء الـ handler التالي
		next.ServeHTTP(w, r)
	})
}
