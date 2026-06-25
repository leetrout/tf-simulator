package sim

import (
	"net/http"

	"github.com/leetrout/terraform-sim/internal/httpx"
)

// Handler serves the /api/sim/* surface backed by a Simulator.
type Handler struct {
	Sim *Simulator
}

// NewHandler returns a sim HTTP handler.
func NewHandler(s *Simulator) *Handler { return &Handler{Sim: s} }

// Register mounts the sim routes on mux.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/sim/snapshot", h.snapshot)
	mux.HandleFunc("GET /api/sim/graph", h.graph)
	mux.HandleFunc("GET /api/sim/status", h.status)
	mux.HandleFunc("GET /api/sim/settings", h.settings)
	mux.HandleFunc("GET /api/sim/scenarios", h.scenarios)
	mux.HandleFunc("POST /api/sim/scenarios/{id}/load", h.loadScenario)
	mux.HandleFunc("POST /api/sim/reset", h.reset)

	mux.HandleFunc("OPTIONS /api/sim/", func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJSON(w, http.StatusNoContent, nil)
	})
}

func (h *Handler) snapshot(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, h.Sim.Snapshot())
}

func (h *Handler) graph(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, h.Sim.Graph(r.URL.Query().Get("view")))
}

func (h *Handler) status(w http.ResponseWriter, r *http.Request) {
	statuses, counts := h.Sim.Status()
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"status": statuses, "counts": counts})
}

func (h *Handler) settings(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, h.Sim.Settings())
}

func (h *Handler) scenarios(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"scenarios": h.Sim.Scenarios()})
}

func (h *Handler) loadScenario(w http.ResponseWriter, r *http.Request) {
	res, err := h.Sim.LoadScenario(r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, res)
}

func (h *Handler) reset(w http.ResponseWriter, r *http.Request) {
	if err := h.Sim.Reset(); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
