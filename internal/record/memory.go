package record

import (
	"context"
	"errors"
	"sort"
)

// ErrNotFound is returned when a record does not exist or is not published.
var ErrNotFound = errors.New("record not found")

// MemoryStore is an in-memory Store for tests and local development.
type MemoryStore struct {
	records map[string]Record
}

// NewMemoryStore returns a store containing the given records keyed by ID.
func NewMemoryStore(records []Record) *MemoryStore {
	recordsByID := make(map[string]Record, len(records))
	for _, r := range records {
		recordsByID[r.ID] = r
	}
	return &MemoryStore{records: recordsByID}
}

// ListPublished returns published records sorted by ID.
func (s *MemoryStore) ListPublished(_ context.Context) ([]Record, error) {
	published := make([]Record, 0)
	for _, r := range s.records {
		if r.Published {
			published = append(published, r)
		}
	}
	sort.Slice(published, func(i, j int) bool {
		return published[i].ID < published[j].ID
	})
	return published, nil
}

// GetPublishedByID returns a published record by ID.
func (s *MemoryStore) GetPublishedByID(_ context.Context, id string) (Record, error) {
	r, ok := s.records[id]
	if !ok || !r.Published {
		return Record{}, ErrNotFound
	}
	return r, nil
}
