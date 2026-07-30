package research

type Asset struct {
	SchemaVersion string         `json:"schema_version"`
	ID            string         `json:"id"`
	Class         string         `json:"class"`
	LocalPath     string         `json:"local_path"`
	SHA256        string         `json:"sha256"`
	ROMsha256     string         `json:"rom_sha256"`
	Confidence    Confidence     `json:"confidence"`
	Provenance    map[string]any `json:"provenance"`
	Validation    map[string]any `json:"validation"`
}
