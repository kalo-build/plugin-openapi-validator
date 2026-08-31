package validate

// Config holds configuration for OpenAPI validation.
type Config struct {
	InputDir       string
	OutputDir      string
	SpecFileName   string
	FailOnWarning  bool
	ReportFileName string
}

// Resolve applies default values for unset fields.
func (c *Config) Resolve() {
	if c.SpecFileName == "" {
		c.SpecFileName = "openapi.yaml"
	}
	if c.ReportFileName == "" {
		c.ReportFileName = "validation-report.json"
	}
}
