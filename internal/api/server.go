package api

import (
	"net/http"
	"strconv"

	"github.com/ThEditor/clutter-paper/internal/api/common"
	"github.com/ThEditor/clutter-paper/internal/api/routes"
	"github.com/ThEditor/clutter-paper/internal/log"
	"github.com/ThEditor/clutter-paper/internal/storage"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "false")

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Start(address string, port int, clickhouse *storage.ClickHouseStorage, redis *storage.RedisStorage, postgres *storage.PostgresStorage) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello, World!"))
	})

	mux.HandleFunc("/api/event", func(w http.ResponseWriter, r *http.Request) {
		routes.PostEvent(w, r, &common.Server{
			Clickhouse: clickhouse,
			Redis:      redis,
			Postgres:   postgres,
		})
	})

	log.Info("API server listening on " + address + ":" + strconv.Itoa(port))
	err := http.ListenAndServe(address+":"+strconv.Itoa(port), corsMiddleware(mux))
	if err != nil {
		log.Info("Server failed to start: " + err.Error())
	}
}
