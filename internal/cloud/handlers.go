package cloud

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/leetrout/terraform-sim/internal/httpx"
)

// Handler serves the /api/cloud/* REST surface backed by a Store.
type Handler struct {
	Store *Store
}

// NewHandler returns a cloud HTTP handler.
func NewHandler(s *Store) *Handler { return &Handler{Store: s} }

// Register mounts the cloud routes on mux using Go's method+wildcard patterns.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/cloud/{collection}", h.list)
	mux.HandleFunc("POST /api/cloud/{collection}", h.create)
	mux.HandleFunc("GET /api/cloud/{collection}/{id}", h.read)
	mux.HandleFunc("PATCH /api/cloud/{collection}/{id}", h.patch)
	mux.HandleFunc("PUT /api/cloud/{collection}/{id}", h.replace)
	mux.HandleFunc("DELETE /api/cloud/{collection}/{id}", h.delete)
	mux.HandleFunc("OPTIONS /api/cloud/{collection}", cors)
	mux.HandleFunc("OPTIONS /api/cloud/{collection}/{id}", cors)
}

func cors(w http.ResponseWriter, r *http.Request) { httpx.WriteJSON(w, http.StatusNoContent, nil) }

func (h *Handler) collectionOr404(w http.ResponseWriter, r *http.Request) (Kind, bool) {
	col := r.PathValue("collection")
	k, ok := KindByCollection(col)
	if !ok {
		httpx.WriteError(w, http.StatusNotFound, "unknown collection: "+col)
		return Kind{}, false
	}
	return k, true
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	k, ok := h.collectionOr404(w, r)
	if !ok {
		return
	}
	items, _ := h.Store.List(k.Collection)
	httpx.WriteJSON(w, http.StatusOK, items)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	k, ok := h.collectionOr404(w, r)
	if !ok {
		return
	}
	obj, _ := New(k.Collection)
	if err := json.NewDecoder(r.Body).Decode(obj); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid body: "+err.Error())
		return
	}
	created, err := h.Store.Create(obj)
	if err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, created)
}

func (h *Handler) read(w http.ResponseWriter, r *http.Request) {
	k, ok := h.collectionOr404(w, r)
	if !ok {
		return
	}
	obj, found := h.Store.Get(k.Collection, r.PathValue("id"))
	if !found {
		httpx.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, obj)
}

func (h *Handler) patch(w http.ResponseWriter, r *http.Request) {
	k, ok := h.collectionOr404(w, r)
	if !ok {
		return
	}
	var fields map[string]any
	if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid body: "+err.Error())
		return
	}
	updated, err := h.Store.Patch(k.Collection, r.PathValue("id"), fields)
	if err != nil {
		h.writeMutationErr(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, updated)
}

func (h *Handler) replace(w http.ResponseWriter, r *http.Request) {
	k, ok := h.collectionOr404(w, r)
	if !ok {
		return
	}
	obj, _ := New(k.Collection)
	if err := json.NewDecoder(r.Body).Decode(obj); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid body: "+err.Error())
		return
	}
	id := r.PathValue("id")
	if _, found := h.Store.Get(k.Collection, id); !found {
		httpx.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	updated, err := h.Store.Put(k.Collection, id, obj)
	if err != nil {
		h.writeMutationErr(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, updated)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	k, ok := h.collectionOr404(w, r)
	if !ok {
		return
	}
	if err := h.Store.Delete(k.Collection, r.PathValue("id")); err != nil {
		h.writeMutationErr(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"deleted": r.PathValue("id")})
}

func (h *Handler) writeMutationErr(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrNotFound) {
		httpx.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
}
