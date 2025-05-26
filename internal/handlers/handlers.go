package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gitvam/platform-go-challenge/internal/models"
	"github.com/gitvam/platform-go-challenge/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/gitvam/platform-go-challenge/internal/middleware"
	"github.com/gitvam/platform-go-challenge/internal/utils"
)

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
// @Success      200 {object} utils.SuccessResponse
// @Failure      401 {object} utils.ErrorResponse "Unauthorized - missing or invalid token"
// @Failure      500 {object} utils.ErrorResponse "Internal server error"
// @Router       /v1/users/{userID}/favorites [get]
func (h *Handler) ListFavorites(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDOrAbort(w, r)
	if !ok {
		return
	}
	favorites, err := h.Store.ListFavorites(userID)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(utils.SuccessResponse{
		Status: "success",
		Data:   favorites,
	})
}

// AddFavorite godoc
// @Summary      Add a favorite asset
// @Description  Add a new favorite asset for the user (chart, insight, or audience).
// @Tags         favorites
// @Param        userID path string true "User ID"
// @Param        asset body models.Asset true "Asset to add"
// @Success      201 {object} utils.SuccessResponse
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} utils.ErrorResponse
// @Failure      409 {object} utils.ErrorResponse
// @Router       /v1/users/{userID}/favorites [post]
func (h *Handler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDOrAbort(w, r)
	if !ok {
		return
	}

	var raw map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		utils.WriteJSONError(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	asset, err := models.DecodeAsset(raw)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Store.AddFavorite(userID, asset); err != nil {
		if err.Error() == "asset already in favorites" {
			utils.WriteJSONError(w, err.Error(), http.StatusConflict)
			return
		}
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(utils.SuccessResponse{
		Status: "success",
		Data:   asset,
	})
}

// RemoveFavorite godoc
// @Summary      Remove a favorite asset
// @Description  Remove an asset from the user's favorites by asset ID.
// @Tags         favorites
// @Param        userID path string true "User ID"
// @Param        assetID path string true "Asset ID"
// @Success      204 "No Content"
// @Failure      401 {object} utils.ErrorResponse
// @Failure      404 {object} utils.ErrorResponse
// @Router       /v1/users/{userID}/favorites/{assetID} [delete]
func (h *Handler) RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDOrAbort(w, r)
	if !ok {
		return
	}
	assetID := chi.URLParam(r, "assetID")
	if err := h.Store.RemoveFavorite(userID, assetID); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusNotFound)
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
// @Success      200 {object} utils.SuccessResponse
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} utils.ErrorResponse
// @Failure      404 {object} utils.ErrorResponse
// @Router       /v1/users/{userID}/favorites/{assetID} [patch]
func (h *Handler) EditFavoriteDescription(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDOrAbort(w, r)
	if !ok {
		return
	}
	assetID := chi.URLParam(r, "assetID")
	var req struct {
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if err := h.Store.EditFavoriteDescription(userID, assetID, req.Description); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(utils.SuccessResponse{
		Status: "success",
		Data: map[string]string{
			"asset_id":    assetID,
			"description": req.Description,
		},
	})
}

func getUserIDOrAbort(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		utils.WriteJSONError(w, "user ID missing from context", http.StatusUnauthorized)
	}
	return userID, ok
}
