package schedule

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h Handler) GetFields(w http.ResponseWriter, r *http.Request) {
	f, err := h.service.GetFields(r.Context())
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{Error: err.Error()})
		return
	}

	render.JSON(w, r, f)
}

func (h Handler) GetGroupsFromID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	g, err := h.service.GetGroups(r.Context(), id)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{Error: err.Error()})
		return
	}

	render.JSON(w, r, g)
}

func (h Handler) GetScheduleFromID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	f := Filter{}

	if d := r.URL.Query().Get("day"); d != "" {
		f.Day = &d
	} else if w := r.URL.Query().Get("week"); w != "" {
		f.Week = &w
	}

	if sg := r.URL.Query().Get("subgroup"); sg != "" {
		f.Subgroup = ParseSubgroup(sg)
	}

	s, err := h.service.GetFilteredSchedule(r.Context(), id, f)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{Error: err.Error()})
		return
	}

	render.JSON(w, r, s)
}
