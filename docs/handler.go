package docs

import (
	_ "embed"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

//go:embed openapi.yaml
var spec []byte

func Spec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	render.Status(r, http.StatusOK)
	if _, err := w.Write(spec); err != nil {
		slog.Error("failed to write openapi spec", "err", err)
	}
}
