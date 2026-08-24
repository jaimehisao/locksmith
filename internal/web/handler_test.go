package web_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jaimehisao/locksmith/internal/copy"
	"github.com/jaimehisao/locksmith/internal/record"
	"github.com/jaimehisao/locksmith/internal/web"
)

func TestUnpublishedRecordReturns404(t *testing.T) {
	store := record.NewMemoryStore([]record.Record{
		{
			ID:         "hidden-1",
			Name:       "Hidden Locksmith",
			Phone:      "415-555-0100",
			License:    "LIC-HIDDEN",
			SourceURL:  "https://example.com/hidden",
			Confidence: 2,
			Published:  false,
		},
	})

	mux := web.NewMux(store)
	req := httptest.NewRequest(http.MethodGet, "/records/hidden-1", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestPublishedDetailShowsRequiredFields(t *testing.T) {
	store := record.NewMemoryStore([]record.Record{
		{
			ID:         "pub-1",
			Name:       "Bay Area Locksmith",
			Phone:      "415-555-0199",
			License:    "LIC-12345",
			SourceURL:  "https://example.com/source",
			Confidence: 4,
			Published:  true,
		},
	})

	mux := web.NewMux(store)
	req := httptest.NewRequest(http.MethodGet, "/records/pub-1", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	for _, want := range []string{
		"Bay Area Locksmith",
		"415-555-0199",
		"LIC-12345",
		"https://example.com/source",
		"4 " + copy.ReportsIndicated,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("response missing %q", want)
		}
	}
}

func TestPublishedListShowsRequiredFields(t *testing.T) {
	store := record.NewMemoryStore([]record.Record{
		{
			ID:         "pub-1",
			Name:       "Bay Area Locksmith",
			Phone:      "415-555-0199",
			License:    "LIC-12345",
			SourceURL:  "https://example.com/source",
			Confidence: 4,
			Published:  true,
		},
		{
			ID:         "hidden-1",
			Name:       "Hidden Locksmith",
			Phone:      "415-555-0100",
			License:    "LIC-HIDDEN",
			SourceURL:  "https://example.com/hidden",
			Confidence: 1,
			Published:  false,
		},
	})

	mux := web.NewMux(store)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	for _, want := range []string{
		"Bay Area Locksmith",
		"415-555-0199",
		"LIC-12345",
		"https://example.com/source",
		"4 " + copy.ReportsIndicated,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("response missing %q", want)
		}
	}
	if strings.Contains(body, "Hidden Locksmith") {
		t.Error("unpublished record appeared in list response")
	}
}

func TestListEmptyStateUsesCopy(t *testing.T) {
	fixedDate := time.Date(2026, 8, 24, 12, 0, 0, 0, time.UTC)
	store := record.NewMemoryStore(nil)
	mux := web.NewMux(store, web.WithNow(func() time.Time { return fixedDate }))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	want := strings.Replace(copy.NoRecordFoundOnDate, "DATE", "2026-08-24", 1)
	if !strings.Contains(rec.Body.String(), want) {
		t.Fatalf("response missing empty-state copy %q; body=%q", want, rec.Body.String())
	}
}
