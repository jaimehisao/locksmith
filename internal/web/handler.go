package web

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/jaimehisao/locksmith/internal/copy"
	"github.com/jaimehisao/locksmith/internal/record"
)

// Option configures a Handler.
type Option func(*Handler)

// WithNow sets the clock used for empty-state copy. Intended for tests.
func WithNow(now func() time.Time) Option {
	return func(h *Handler) {
		h.Now = now
	}
}

// Handler serves published list and detail pages.
type Handler struct {
	Store Store
	Now   func() time.Time
}

// Store is the subset of record.Store used by HTTP handlers.
type Store interface {
	ListPublished(ctx context.Context) ([]record.Record, error)
	GetPublishedByID(ctx context.Context, id string) (record.Record, error)
}

type listPageData struct {
	Records    []record.Record
	EmptyState string
}

type detailPageData struct {
	Record         record.Record
	ConfidenceText string
	Related        []record.Record
	RelatedLabel   string
}

var (
	listTemplate   *template.Template
	detailTemplate *template.Template
)

func init() {
	funcMap := template.FuncMap{
		"confidence": confidenceText,
	}
	listTemplate = template.Must(template.New("list").Funcs(funcMap).Parse(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="utf-8"><title>LocksmithWatch SF</title></head>
<body>
{{if .EmptyState}}
<p>{{.EmptyState}}</p>
{{else}}
<ul>
{{range .Records}}
<li><a href="/records/{{.ID}}">{{.Name}}</a> {{.Phone}} {{.License}} <a href="{{.SourceURL}}">{{.SourceURL}}</a> {{confidence .Confidence}}</li>
{{end}}
</ul>
{{end}}
</body>
</html>`))

	detailTemplate = template.Must(template.New("detail").Parse(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="utf-8"><title>{{.Record.Name}}</title></head>
<body>
<p>{{.Record.Name}}</p>
<p>{{.Record.Phone}}</p>
<p>{{.Record.License}}</p>
<p><a href="{{.Record.SourceURL}}">{{.Record.SourceURL}}</a></p>
<p>{{.ConfidenceText}}</p>
{{if .Related}}
<p>{{.RelatedLabel}}</p>
<ul>
{{range .Related}}
<li><a href="/records/{{.ID}}">{{.Name}}</a> {{.Phone}} {{.License}}</li>
{{end}}
</ul>
{{end}}
<p><a href="/">Back</a></p>
</body>
</html>`))
}

// NewMux returns an HTTP mux for published list and detail routes.
func NewMux(store Store, opts ...Option) http.Handler {
	h := &Handler{Store: store, Now: time.Now}
	for _, opt := range opts {
		opt(h)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", h.serveList)
	mux.HandleFunc("GET /records/{id}", h.serveDetail)
	return mux
}

func (h *Handler) serveList(w http.ResponseWriter, r *http.Request) {
	records, err := h.Store.ListPublished(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data := listPageData{Records: records}
	if len(records) == 0 {
		data.EmptyState = noRecordFoundOn(h.now())
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := listTemplate.Execute(w, data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *Handler) serveDetail(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rec, err := h.Store.GetPublishedByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, record.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	related, err := h.relatedPublished(r.Context(), rec)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data := detailPageData{
		Record:         rec,
		ConfidenceText: confidenceText(rec.Confidence),
		Related:        related,
		RelatedLabel:   copy.AppearRelated,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := detailTemplate.Execute(w, data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *Handler) relatedPublished(ctx context.Context, current record.Record) ([]record.Record, error) {
	all, err := h.Store.ListPublished(ctx)
	if err != nil {
		return nil, err
	}

	related := make([]record.Record, 0)
	for _, candidate := range all {
		if candidate.ID == current.ID {
			continue
		}
		if candidate.License == current.License || candidate.Phone == current.Phone {
			related = append(related, candidate)
		}
	}
	return related, nil
}

func (h *Handler) now() time.Time {
	if h.Now != nil {
		return h.Now()
	}
	return time.Now()
}
