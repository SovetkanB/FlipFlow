package expense

import (
	"errors"
	"net/http"

	"github.com/SovetkanB/FlipFlow/internal/domain/auth"
	"github.com/SovetkanB/FlipFlow/internal/pkg/response"
	"github.com/SovetkanB/FlipFlow/internal/pkg/validator"
	"github.com/go-chi/chi"
)

type Handler struct {
	service Service
}

func NewHandler(srv Service) *Handler {
	return &Handler{service: srv}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	if claims == nil {
		response.Unauthorized(w)
		return
	}

	var req CreateExpenseRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	projectID := chi.URLParam(r, "projectID")

	res, err := h.service.Create(r.Context(), projectID, claims.UserID, req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	if claims == nil {
		response.Unauthorized(w)
		return
	}

	projectID := chi.URLParam(r, "projectID")
	res, err := h.service.List(r.Context(), projectID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, res)

}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	if claims == nil {
		response.Unauthorized(w)
		return
	}

	projectID := chi.URLParam(r, "projectID")
	expenseID := chi.URLParam(r, "expenseID")

	err := h.service.Delete(r.Context(), projectID, expenseID)
	if err != nil {
		switch {
		case errors.Is(err, response.ErrNotFound):
			response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		default:
			response.Error(w, http.StatusInternalServerError, "Error", err.Error())
			return
		}
	}
}
