package schedule

import (
	"net/http"
	"uz-plan-api/internal/errs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"golang.org/x/time/rate"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type Handler struct {
	service *Service
	limiter *rate.Limiter
}

func NewHandler(service *Service, limiter *rate.Limiter) *Handler {
	return &Handler{service: service, limiter: limiter}
}

func (h Handler) GetFields(w http.ResponseWriter, r *http.Request) {
	if !h.checkRateLimit(w, r) {
		return
	}

	f, err := h.service.GetFields(r.Context())
	if err != nil {
		render.Status(r, errs.StatusFromErr(err))
		render.JSON(w, r, ErrorResponse{Error: err.Error()})
		return
	}

	render.JSON(w, r, f)
}

func (h Handler) GetGroupsFromID(w http.ResponseWriter, r *http.Request) {
	if !h.checkRateLimit(w, r) {
		return
	}

	id := chi.URLParam(r, "id")
	if !h.isID(w, r, id) {
		return
	}

	g, err := h.service.GetGroups(r.Context(), id)
	if err != nil {
		render.Status(r, errs.StatusFromErr(err))
		render.JSON(w, r, ErrorResponse{Error: err.Error()})
		return
	}

	render.JSON(w, r, g)
}

func (h Handler) GetScheduleFromID(w http.ResponseWriter, r *http.Request) {
	if !h.checkRateLimit(w, r) {
		return
	}

	id := chi.URLParam(r, "id")
	if !h.isID(w, r, id) {
		return
	}

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
		render.Status(r, errs.StatusFromErr(err))
		render.JSON(w, r, ErrorResponse{Error: err.Error()})
		return
	}

	render.JSON(w, r, s)
}

func (h Handler) checkRateLimit(w http.ResponseWriter, r *http.Request) bool {
	if h.limiter.Allow() {
		return true
	}
	render.Status(r, http.StatusTooManyRequests)
	render.JSON(w, r, ErrorResponse{Error: errs.ErrTooManyReq.Error()})
	return false
}

func (h Handler) isID(w http.ResponseWriter, r *http.Request, id string) bool {
	if id == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{Error: "id is required"})
		return false
	}

	return true
}
