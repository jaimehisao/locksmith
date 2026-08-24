package record

import "context"

// Store provides access to locksmith records.
type Store interface {
	ListPublished(ctx context.Context) ([]Record, error)
	GetPublishedByID(ctx context.Context, id string) (Record, error)
}
