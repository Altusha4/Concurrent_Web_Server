package handler

import (
	"assignment2/internal/models"
	"assignment2/internal/service"
	"encoding/json"
	"net/http"
)

type DataHandler struct {
	service *service.DataService
}

func NewDataHandler(service *service.DataService) *DataHandler {
	return &DataHandler{service: service}
}

func (h *DataHandler) PostData(w http.ResponseWriter, r *http.Request) {
	var kv models.KeyValue

	if err := json.NewDecoder(r.Body).Decode(&kv); err != nil {
		h.sendError(w, "Invalid JSON. Expected: {\"key\": \"value\"}", http.StatusBadRequest)
		return
	}

	if kv.Key == "" || kv.Value == "" {
		h.sendError(w, "Both 'key' and 'value' are required", http.StatusBadRequest)
		return
	}

	h.service.SaveKeyValue(kv.Key, kv.Value)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "created",
		"key":    kv.Key,
	})
}

func (h *DataHandler) GetData(w http.ResponseWriter, r *http.Request) {
	data := h.service.GetAll()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *DataHandler) DeleteData(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if key == "" {
		h.sendError(w, "Key is required", http.StatusBadRequest)
		return
	}

	if !h.service.DeleteKey(key) {
		h.sendError(w, "Key not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "deleted",
		"key":    key,
	})
}

func (h *DataHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	totalRequests, dbSize := h.service.GetStats()

	stats := map[string]interface{}{
		"total_requests": totalRequests,
		"database_size":  dbSize,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *DataHandler) sendError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
