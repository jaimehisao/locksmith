package locksmithwatch

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jaimehisao/locksmith/internal/copy"
)

func samplePublishedRecord() Record {
	return Record{
		ID:         "alpha-locksmith",
		Name:       "Alpha Locksmith",
		Phone:      "415-555-0100",
		License:    "LCO-1234",
		SourceURL:  "https://example.com/alpha",
		Confidence: "high",
		Published:  true,
	}
}

func sampleUnpublishedRecord() Record {
	record := samplePublishedRecord()
	record.ID = "hidden-locksmith"
	record.Name = "Hidden Locksmith"
	record.Published = false
	return record
}

func TestDetail_UnpublishedReturns404(t *testing.T) {
	store := NewMemoryStore(sampleUnpublishedRecord())
	handler := NewHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/locksmiths/hidden-locksmith", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestDetail_MissingReturns404(t *testing.T) {
	store := NewMemoryStore()
	handler := NewHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/locksmiths/missing", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestDetail_PublishedShowsFields(t *testing.T) {
	record := samplePublishedRecord()
	store := NewMemoryStore(record)
	handler := NewHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/locksmiths/alpha-locksmith", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	for _, want := range []string{
		record.Name,
		record.Phone,
		record.License,
		record.SourceURL,
		record.Confidence,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("detail body missing %q", want)
		}
	}
}

func TestList_EmptyStateUsesAllowListedCopy(t *testing.T) {
	fixedNow := time.Date(2026, 8, 24, 12, 0, 0, 0, time.UTC)
	handler := &Handler{
		store: NewMemoryStore(),
		now:   func() time.Time { return fixedNow },
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	for _, phrase := range []string{
		copy.ReportsIndicated,
		copy.AppearRelated,
		EmptyStateNoRecordFound(fixedNow),
	} {
		if !strings.Contains(body, phrase) {
			t.Errorf("empty list body missing allow-listed phrase %q", phrase)
		}
	}
}

func TestList_PublishedRecordsOnly(t *testing.T) {
	published := samplePublishedRecord()
	unpublished := sampleUnpublishedRecord()
	handler := NewHandler(NewMemoryStore(published, unpublished))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	if !strings.Contains(body, published.Name) {
		t.Errorf("list body missing published record %q", published.Name)
	}
	if strings.Contains(body, unpublished.Name) {
		t.Errorf("list body unexpectedly contains unpublished record %q", unpublished.Name)
	}
}

func TestMemoryStore_ListPublishedSortedByName(t *testing.T) {
	first := samplePublishedRecord()
	first.ID = "b-record"
	first.Name = "Bravo Locksmith"

	second := samplePublishedRecord()
	second.ID = "a-record"
	second.Name = "Alpha Locksmith"

	store := NewMemoryStore(first, second)
	got := store.ListPublished()
	if len(got) != 2 {
		t.Fatalf("published count = %d, want 2", len(got))
	}
	if got[0].Name != "Alpha Locksmith" || got[1].Name != "Bravo Locksmith" {
		t.Fatalf("published order = [%q, %q], want [Alpha Locksmith, Bravo Locksmith]", got[0].Name, got[1].Name)
	}
}
