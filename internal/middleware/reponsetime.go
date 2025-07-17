package middleware

import (
	"log"
	"net/http"
	"time"
)

// Middleware لحساب وقت الاستجابة وطباعته في اللوج
func ResponseTimeMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// نغلف ResponseWriter لمتابعة حالة الاستجابة
		wrappedWriter := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		// بعد تنفيذ الطلب: احسب الوقت المستغرق
		duration := time.Since(start)
		wrappedWriter.Header().Set("x-Response-Time", duration.String())

		// استدعاء handler الرئيسي
		next.ServeHTTP(wrappedWriter, r)

		log.Printf("[%s] %s %d %s", r.Method, r.URL.Path, wrappedWriter.status, duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
