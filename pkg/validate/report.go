package validate

// Report holds the results of an OpenAPI validation run.
type Report struct {
	Valid    bool    `json:"valid"`
	Errors   []Issue `json:"errors"`
	Warnings []Issue `json:"warnings"`
	SpecFile string  `json:"specFile"`
	Version  string  `json:"openapiVersion,omitempty"`
}

// Issue represents a single validation error or warning.
type Issue struct {
	Message  string `json:"message"`
	Location string `json:"location,omitempty"`
	Rule     string `json:"rule,omitempty"`
}
