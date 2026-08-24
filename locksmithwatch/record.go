package locksmithwatch

// Record holds locksmith data shown on published pages.
type Record struct {
	ID         string
	Name       string
	Phone      string
	License    string
	SourceURL  string
	Confidence string
	Published  bool
}
