package locksmithwatch

import (
	"net/http"
	"strings"
	"time"
)

// Handler serves published list and detail pages.
type Handler struct {
	store Store
	now   func() time.Time
}

// NewHandler returns HTTP handlers backed by store.
func NewHandler(store Store) http.Handler {
	return (&Handler{
		store: store,
		now:   time.Now,
	}).routes()
}

func (h *Handler) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.serveList)
	mux.HandleFunc("GET /locksmiths/{id}", h.serveDetail)
	return mux
}

func (h *Handler) serveList(w http.ResponseWriter, r *http.Request) {
	body, err := RenderListPage(h.store.ListPublished(), h.now())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(body))
}

func (h *Handler) serveDetail(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	record, ok := h.store.Get(id)
	if !ok || !record.Published {
		http.NotFound(w, r)
		return
	}

	body, err := RenderDetailPage(record)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(body))
}
