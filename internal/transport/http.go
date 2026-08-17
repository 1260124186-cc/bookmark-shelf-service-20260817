package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/domain"
	"github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/service"
)

type Handler struct {
	library *service.Library
}

func NewHandler(library *service.Library) http.Handler {
	handler := &Handler{library: library}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handler.health)
	mux.HandleFunc("POST /collections", handler.createCollection)
	mux.HandleFunc("GET /collections", handler.listCollections)
	mux.HandleFunc("POST /bookmarks", handler.saveBookmark)
	mux.HandleFunc("POST /bookmarks/{id}/move", handler.moveBookmark)
	mux.HandleFunc("POST /bookmarks/{id}/archive", handler.archiveBookmark)
	mux.HandleFunc("GET /report", handler.report)
	return mux
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) createCollection(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name string `json:"name"`
	}
	if !decodeJSON(w, r, &request) {
		return
	}
	collection, err := h.library.CreateCollection(r.Context(), request.Name)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, collection)
}

func (h *Handler) listCollections(w http.ResponseWriter, r *http.Request) {
	collections, err := h.library.ListCollections(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, collections)
}

func (h *Handler) saveBookmark(w http.ResponseWriter, r *http.Request) {
	var request service.BookmarkInput
	if !decodeJSON(w, r, &request) {
		return
	}
	bookmark, err := h.library.SaveBookmark(r.Context(), request)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, bookmark)
}

func (h *Handler) moveBookmark(w http.ResponseWriter, r *http.Request) {
	var request struct {
		CollectionID string `json:"collection_id"`
	}
	if !decodeJSON(w, r, &request) {
		return
	}
	bookmark, err := h.library.MoveBookmark(r.Context(), r.PathValue("id"), request.CollectionID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, bookmark)
}

func (h *Handler) archiveBookmark(w http.ResponseWriter, r *http.Request) {
	bookmark, err := h.library.ArchiveBookmark(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, bookmark)
}

func (h *Handler) report(w http.ResponseWriter, r *http.Request) {
	report, err := h.library.BuildReport(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return false
	}
	return true
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidBookmark), errors.Is(err, domain.ErrInvalidCollection):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, domain.ErrBookmarkNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.Is(err, domain.ErrDuplicateBookmark), errors.Is(err, domain.ErrCollectionNotFound):
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		writeJSON(w, http.StatusRequestTimeout, map[string]string{"error": "request canceled"})
	default:
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func pathID(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}
