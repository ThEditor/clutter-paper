package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/ThEditor/clutter-paper/internal/storage"
)

type RequestData struct {
	VisitorUserAgent string
	SiteID           string
	Referrer         string
	Page             string
}

func validate(data RequestData) error {
	if strings.TrimSpace(data.VisitorUserAgent) == "" {
		return errors.New("visitor user agent is required")
	}
	if strings.TrimSpace(data.SiteID) == "" {
		return errors.New("site ID is required")
	}
	if strings.TrimSpace(data.Page) == "" {
		return errors.New("page is required")
	}
	return nil
}

func PostEvent(w http.ResponseWriter, r *http.Request, clickhouse *storage.ClickHouseStorage) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var data RequestData

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&data); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate(data); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// check for site_id in cache/postgres

	visitorIP := r.Header.Get("X-Forwarded-For")
	if visitorIP == "" {
		visitorIP = r.RemoteAddr
	}

	eventData := storage.EventData{
		VisitorIP:        visitorIP,
		VisitorUserAgent: data.VisitorUserAgent,
		SiteID:           data.SiteID,
		Referrer:         data.Referrer,
		Page:             data.Page,
	}

	if err := clickhouse.InsertEvent(eventData); err != nil {
		http.Error(w, "Failed to store event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Data received successfully",
	})
}
