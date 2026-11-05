package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/service"
)

type PostHandler struct {
	svc service.PostService
}

func NewPostHandler(svc service.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

// CreatePost
// Post /api/post
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var request model.PostCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if request.Title == "" {
		respondWithError(w, http.StatusBadRequest, "title is required")
		return
	}

	post, err := h.svc.Create(request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, post)
}

// List
// GET /api/post/page?offset=0&limit=10
func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	offset, err := parseQueryInt(r, "offset", 0)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid offset")
		return
	}
	limit, err := parseQueryInt(r, "limit", 10)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid limit")
		return
	}

	posts, err := h.svc.List(offset, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, posts)
}

// GetByID
// GET /api/post?id=1
func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	raw := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	post, err := h.svc.GetById(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "post not found")
		return
	}
	respondWithJSON(w, http.StatusOK, post)
}

// parseQueryInt
func parseQueryInt(r *http.Request, key string, fallback int) (int, error) {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, err
	}
	if value < 0 {
		return 0, errors.New("must not be negative")
	}
	return value, nil
}
