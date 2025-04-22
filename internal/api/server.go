package api

import (
	"net/http"
	"strconv"

	"github.com/ThEditor/clutter-paper/internal/api/common"
	"github.com/ThEditor/clutter-paper/internal/api/routes"
	"github.com/ThEditor/clutter-paper/internal/log"
	"github.com/ThEditor/clutter-paper/internal/storage"
)

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
	err := http.ListenAndServe(address+":"+strconv.Itoa(port), mux)
	if err != nil {
		log.Info("Server failed to start: " + err.Error())
	}
}
