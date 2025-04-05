package api

import (
	"net/http"
	"strconv"

	"github.com/ThEditor/clutter-paper/internal/api/routes"
	"github.com/ThEditor/clutter-paper/internal/log"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:5173" || origin == "http://127.0.0.1:5173" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}
		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Start(address string, port int) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello, World!"))
	})

	mux.HandleFunc("/api/event", routes.PostEvent)

	log.Info("API server listening on " + address + ":" + strconv.Itoa(port))
	err := http.ListenAndServe(address+":"+strconv.Itoa(port), corsMiddleware(mux))
	if err != nil {
		log.Info("Server failed to start: " + err.Error())
	}
}
