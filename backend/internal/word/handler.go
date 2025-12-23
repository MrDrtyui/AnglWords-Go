package word

import (
	"app/domain"
	"app/internal/middleware"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
}

// GetWord godoc
// @Summary      Get a word
// @Description  Get a word by its string representation
// @Tags         words
// @Produce      json
// @Param        word path string true "Word to get"
// @Success      200 {object} domain.WordResponse "Successfully retrieved"
// @Failure      404 {object} domain.ErrorResponse "Word not found"
// @Failure      500 {object} domain.ErrorResponse "Internal server error"
// @Router       /word/{word} [get]
func (h *Handler) GetWord(w http.ResponseWriter, r *http.Request) {
	wordStr := chi.URLParam(r, "word")

	word, err := h.Service.GetWord(r.Context(), wordStr)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "word not found")
		return
	}

	respondWithJSON(w, http.StatusOK, ToWordResponse(word))
}

// GetMyWords godoc
// @Summary      Get all my words
// @Description  Get all words for the current user
// @Tags         words
// @Produce      json
// @Success      200 {object} []domain.WordResponse "Successfully retrieved"
// @Failure      401 {object} domain.ErrorResponse "Unauthorized"
// @Failure      500 {object} domain.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /words/my [get]
func (h *Handler) GetMyWords(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	words, err := h.Service.GetMyWords(r.Context(), userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, ToWordResponseSlice(words))
}

// CreateWord godoc
// @Summary      Create a new word
// @Description  Create a new word for the current user
// @Tags         words
// @Accept       json
// @Produce      json
// @Param        request body domain.WordCreateDto true "Word to create"
// @Success      201 {object} domain.WordResponse "Successfully created"
// @Failure      400 {object} domain.ErrorResponse "Invalid request body or validation error"
// @Failure      401 {object} domain.ErrorResponse "Unauthorized"
// @Failure      500 {object} domain.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /words [post]
func (h *Handler) CreateWord(w http.ResponseWriter, r *http.Request) {
	var dto domain.WordCreateDto
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userId, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	word, err := h.Service.CreateWord(r.Context(), dto, userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, ToWordResponse(word))
}

// GetAllWords godoc
// @Summary      Get all words (public)
// @Description  Get all words in the database
// @Tags         words
// @Produce      json
// @Success      200 {object} []domain.WordResponse "Successfully retrieved"
// @Failure      500 {object} domain.ErrorResponse "Internal server error"
// @Router       /words/all [get]
func (h *Handler) GetAllWords(w http.ResponseWriter, r *http.Request) {
	words, err := h.Service.GetAllWords(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, ToWordResponseSlice(words))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
