// Это сервис - заглушка для проверки корректности отправки данных

package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func loggingMiddleware(logger *log.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)
		logger.Printf(
			"[%s] HTTP: %s %s -> %d (%dms)",
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			lrw.statusCode,
			duration.Milliseconds(),
		)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	logFile, err := os.OpenFile("log1.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Ошибка открытия файла логов: %v", err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", 0)

	// Обработчик с middleware для логирования
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		codes := []int{
			http.StatusOK,
			http.StatusCreated,
			http.StatusBadRequest,
			http.StatusUnauthorized,
			http.StatusForbidden,
			http.StatusNotFound,
			http.StatusInternalServerError,
		}
		status := codes[rand.Intn(len(codes))]
		w.WriteHeader(status)
	}

	http.HandleFunc("/users", loggingMiddleware(logger, handler))
	http.ListenAndServe(":8081", nil)
}
