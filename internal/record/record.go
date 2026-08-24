package record

// Record is a locksmith monitoring entry.
type Record struct {
	ID         string
	Name       string
	Phone      string
	License    string
	SourceURL  string
	Confidence int
	Published  bool
}
