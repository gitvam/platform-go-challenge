package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gitvam/platform-go-challenge/internal/models"
	"github.com/gitvam/platform-go-challenge/internal/store"
	"github.com/go-chi/chi/v5"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSONError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: msg})
}

// Handler holds dependencies (store)
type Handler struct {
	Store store.Store
}

// NewHandler creates a new Handler with dependencies injected
func NewHandler(s store.Store) *Handler {
	return &Handler{Store: s}
}

// List all favorites for a user
func (h *Handler) ListFavorites(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	favorites, err := h.Store.ListFavorites(userID)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(favorites)
}

// Add a new favorite for a user
func (h *Handler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	var raw map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		writeJSONError(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	assetType, ok := raw["type"].(string)
	if !ok {
		writeJSONError(w, "missing asset type", http.StatusBadRequest)
		return
	}

	var asset models.Asset
	switch assetType {
	case string(models.AssetTypeChart):
		var chart models.Chart
		data, _ := json.Marshal(raw)
		if err := json.Unmarshal(data, &chart); err != nil {
			writeJSONError(w, "invalid chart", http.StatusBadRequest)
			return
		}
		asset = &chart
	case string(models.AssetTypeInsight):
		var insight models.Insight
		data, _ := json.Marshal(raw)
		if err := json.Unmarshal(data, &insight); err != nil {
			writeJSONError(w, "invalid insight", http.StatusBadRequest)
			return
		}
		asset = &insight
	case string(models.AssetTypeAudience):
		var audience models.Audience
		data, _ := json.Marshal(raw)
		if err := json.Unmarshal(data, &audience); err != nil {
			writeJSONError(w, "invalid audience", http.StatusBadRequest)
			return
		}
		asset = &audience
	default:
		writeJSONError(w, "unknown asset type", http.StatusBadRequest)
		return
	}

	if err := h.Store.AddFavorite(userID, asset); err != nil {
		if err.Error() == "asset already in favorites" {
			writeJSONError(w, err.Error(), http.StatusConflict)
			return
		}
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Remove a favorite
func (h *Handler) RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	assetID := chi.URLParam(r, "assetID")
	if err := h.Store.RemoveFavorite(userID, assetID); err != nil {
		writeJSONError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Edit description of a favorite
func (h *Handler) EditFavoriteDescription(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	assetID := chi.URLParam(r, "assetID")
	var req struct {
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if err := h.Store.EditFavoriteDescription(userID, assetID, req.Description); err != nil {
		writeJSONError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
