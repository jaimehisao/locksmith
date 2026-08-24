package locksmithwatch

import "sort"

// Store provides read access to locksmith records.
type Store interface {
	ListPublished() []Record
	Get(id string) (Record, bool)
}

// MemoryStore keeps records in memory for handlers and tests.
type MemoryStore struct {
	records map[string]Record
}

// NewMemoryStore returns a store seeded with the given records.
func NewMemoryStore(records ...Record) *MemoryStore {
	recordsByID := make(map[string]Record, len(records))
	for _, record := range records {
		recordsByID[record.ID] = record
	}
	return &MemoryStore{records: recordsByID}
}

// ListPublished returns published records sorted by name.
func (s *MemoryStore) ListPublished() []Record {
	published := make([]Record, 0)
	for _, record := range s.records {
		if record.Published {
			published = append(published, record)
		}
	}
	sort.Slice(published, func(i, j int) bool {
		return published[i].Name < published[j].Name
	})
	return published
}

// Get returns a record by ID.
func (s *MemoryStore) Get(id string) (Record, bool) {
	record, ok := s.records[id]
	return record, ok
}
