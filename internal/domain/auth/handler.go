package auth

import (
	"net/http"

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
		response.BadRequest(w, "validation error")
	}

	resp, err := h.service.Register(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error", err.Error())
	}

	response.JSON(w, http.StatusOK, resp)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

}
