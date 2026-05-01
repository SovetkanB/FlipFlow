package auth

import (
	"errors"
	"net/http"

	"github.com/SovetkanB/FlipFlow/internal/api/middleware"
	"github.com/SovetkanB/FlipFlow/internal/pkg/response"
	"github.com/SovetkanB/FlipFlow/internal/pkg/validator"
)

type Handler struct {
	service Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{
		service: svc,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	resp, err := h.service.Register(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, response.ErrEmailTaken):
			response.Error(w, http.StatusConflict, "EMAIL_TAKEN", err.Error())
			return
		default:
			response.Error(w, http.StatusInternalServerError, "Error", err.Error())
			return
		}
	}

	response.JSON(w, http.StatusOK, resp)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	resp, err := h.service.Login(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, resp)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	resp, err := h.service.RefreshTokens(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "INVALID_REFRESH_TOKEN", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, resp)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.ClaimsFromContext(r.Context())
	if claims == nil {
		response.Unauthorized(w)
		return
	}
}
