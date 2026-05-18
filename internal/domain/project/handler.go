package project

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/SovetkanB/FlipFlow/internal/domain/auth"
	"github.com/SovetkanB/FlipFlow/internal/pkg/response"
	"github.com/SovetkanB/FlipFlow/internal/pkg/validator"
	"github.com/go-chi/chi"
)

type Handler struct {
	service *Service
}

func NewHandler(srv *Service) *Handler {
	return &Handler{
		service: srv,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	if claims == nil {
		response.Unauthorized(w)
		return
	}

	f := ListFilter{
		Status: r.URL.Query().Get("status"),
		Limit:  parseInt(r, "limit", 20),
		Offset: parseInt(r, "offset", 0),
	}

	list, total, err := h.service.List(r.Context(), claims.UserID, f)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"items":  list,
		"total":  total,
		"limit":  f.Limit,
		"offset": f.Offset,
	})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	if claims == nil {
		response.Unauthorized(w)
		return
	}

	var req CreateProjectRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	res, err := h.service.Create(r.Context(), claims.UserID, req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	if claims == nil {
		response.Unauthorized(w)
		return
	}

	projectID := chi.URLParam(r, "projectID")

	res, err := h.service.GetByID(r.Context(), projectID, claims.UserID)
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

	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	if claims == nil {
		response.Unauthorized(w)
		return
	}

	projectID := chi.URLParam(r, "projectID")

	var req UpdateProjectRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	res, err := h.service.Update(r.Context(), projectID, claims.UserID, req)
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

	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	if claims == nil {
		response.Unauthorized(w)
		return
	}
	projectID := chi.URLParam(r, "projectID")

	var body struct {
		Status    Status   `json:"status"`
		SoldPrice *float64 `json:"sold_price"`
	}
	if err := validator.DecodeAndValidate(r, &body); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	res, err := h.service.ChangeStatus(r.Context(), projectID, claims.UserID, body.Status, body.SoldPrice)
	if err != nil {
		response.BadRequest(w, err.Error())
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

	err := h.service.Delete(r.Context(), projectID, claims.UserID)
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

func (h *Handler) FinancialSummary(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromContext(r.Context())
	if claims == nil {
		response.Unauthorized(w)
		return
	}

	projectID := chi.URLParam(r, "projectID")
	res, err := h.service.GetFinancialSummary(r.Context(), projectID, claims.UserID)
	if err != nil {
		response.NotFound(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, res)
}

func parseInt(r *http.Request, key string, def int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
