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

// ListFavorites godoc
// @Summary      List all favorites for a user
// @Description  Get all favorite assets (charts, insights, audiences) for the specified user.
// @Tags         favorites
// @Param        userID path string true "User ID"
// @Success      200 {array} models.Asset "List of favorite assets"
// @Failure      401 {object} handlers.ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      500 {object} handlers.ErrorResponse "Internal server error"
// @Router       /v1/users/{userID}/favorites [get]
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

// AddFavorite godoc
// @Summary      Add a favorite asset
// @Description  Add a new favorite asset for the user (chart, insight, or audience). Provide "type" field as 'chart', 'insight', or 'audience'.
// @Tags         favorites
// @Param        userID path string true "User ID"
// @Param        asset body models.Asset true "Asset to add. Set the 'type' field to 'chart', 'insight', or 'audience'."
// @Success      201 {object} models.Asset "Favorite asset created"
// @Failure      400 {object} handlers.ErrorResponse "Bad request - invalid input"
// @Failure      401 {object} handlers.ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      409 {object} handlers.ErrorResponse "Conflict - asset already in favorites"
// @Router       /v1/users/{userID}/favorites [post]
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(asset)
}

// RemoveFavorite godoc
// @Summary      Remove a favorite asset
// @Description  Remove an asset from the user's favorites by asset ID.
// @Tags         favorites
// @Param        userID path string true "User ID"
// @Param        assetID path string true "Asset ID"
// @Success      204 "No Content - asset deleted"
// @Failure      401 {object} handlers.ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      404 {object} handlers.ErrorResponse "Not found - asset not found"
// @Router       /v1/users/{userID}/favorites/{assetID} [delete]
func (h *Handler) RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	assetID := chi.URLParam(r, "assetID")
	if err := h.Store.RemoveFavorite(userID, assetID); err != nil {
		writeJSONError(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// EditFavoriteDescription godoc
// @Summary      Edit favorite asset description
// @Description  Edit the description of a favorite asset.
// @Tags         favorites
// @Param        userID path string true "User ID"
// @Param        assetID path string true "Asset ID"
// @Param        body body handlers.EditDescriptionRequest true "New Description"
// @Success      200 {object} models.Asset "Updated asset"
// @Failure      400 {object} handlers.ErrorResponse "Bad request - invalid input"
// @Failure      401 {object} handlers.ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      404 {object} handlers.ErrorResponse "Not found - asset not found"
// @Router       /v1/users/{userID}/favorites/{assetID} [patch]
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
