package schema

// Provenance records where a recovered value came from.
type Provenance struct {
	Source     string `json:"source"`
	Confidence string `json:"confidence"`
}
