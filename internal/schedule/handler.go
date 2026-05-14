package schedule

import (
	"net/http"
	"regexp"
	"strings"
	"uz-plan-api/internal/errs"
	"uz-plan-api/internal/model"

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
		render.Status(r, errs.StatusFromErr(err))
		render.JSON(w, r, ErrorResponse{Error: err.Error()})
		return
	}

	render.JSON(w, r, f)
}

func (h Handler) GetGroupsFromID(w http.ResponseWriter, r *http.Request) {
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
	id := chi.URLParam(r, "id")
	if !h.isID(w, r, id) {
		return
	}

	f := model.Filter{}

	if d := r.URL.Query().Get("day"); d != "" {
		f.Day = &d
	} else if w := r.URL.Query().Get("week"); w != "" {
		f.Week = &w
	}

	if sg := r.URL.Query().Get("subgroup"); sg != "" {
		sgUpper := strings.ToUpper(sg)
		f.Subgroup = model.ParseSubgroup(sgUpper)
	}

	s, err := h.service.GetFilteredSchedule(r.Context(), id, f)
	if err != nil {
		render.Status(r, errs.StatusFromErr(err))
		render.JSON(w, r, ErrorResponse{Error: err.Error()})
		return
	}

	render.JSON(w, r, s)
}

var validID = regexp.MustCompile(`^-?\d+$`)

func (h Handler) isID(w http.ResponseWriter, r *http.Request, id string) bool {
	if id == "" || !validID.MatchString(id) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{Error: "id is required"})
		return false
	}

	return true
}
