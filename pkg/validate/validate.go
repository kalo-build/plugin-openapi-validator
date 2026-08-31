package validate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pb33f/libopenapi"
	validator "github.com/pb33f/libopenapi-validator"
)

// ValidateSpec reads an OpenAPI spec, validates it using libopenapi, and
// writes a JSON report to the output directory.
func ValidateSpec(cfg Config) (*Report, error) {
	specPath := filepath.Join(cfg.InputDir, cfg.SpecFileName)
	specBytes, err := os.ReadFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec file '%s': %w", specPath, err)
	}

	report := ValidateBytes(specBytes, cfg.SpecFileName)

	if err := writeReport(cfg, report); err != nil {
		return nil, fmt.Errorf("failed to write report: %w", err)
	}

	return report, nil
}

// ValidateBytes validates raw OpenAPI spec bytes and returns a report.
func ValidateBytes(specBytes []byte, fileName string) *Report {
	report := &Report{
		Valid:    true,
		Errors:   []Issue{},
		Warnings: []Issue{},
		SpecFile: fileName,
	}

	document, docErr := libopenapi.NewDocument(specBytes)
	if docErr != nil {
		report.Valid = false
		report.Errors = append(report.Errors, Issue{
			Message: fmt.Sprintf("Failed to parse OpenAPI document: %s", docErr.Error()),
			Rule:    "document-parse",
		})
		return report
	}

	model, modelErr := document.BuildV3Model()
	if modelErr != nil {
		report.Errors = append(report.Errors, Issue{
			Message: modelErr.Error(),
			Rule:    "model-build",
		})
	}

	if model != nil && model.Model.Info != nil {
		report.Version = model.Model.Info.Version
	}

	docValidator, validatorErrs := validator.NewValidator(document)
	if validatorErrs != nil {
		for _, e := range validatorErrs {
			report.Errors = append(report.Errors, Issue{
				Message: e.Error(),
				Rule:    "validator-init",
			})
		}
	}

	if docValidator != nil {
		_, validationErrs := docValidator.ValidateDocument()
		for _, ve := range validationErrs {
			msg := ve.Message
			if ve.Reason != "" {
				msg = fmt.Sprintf("%s: %s", ve.Message, ve.Reason)
			}

			location := ""
			if ve.SchemaValidationErrors != nil {
				var locs []string
				for _, sve := range ve.SchemaValidationErrors {
					if sve.Location != "" {
						locs = append(locs, sve.Location)
					}
				}
				if len(locs) > 0 {
					location = strings.Join(locs, ", ")
				}
			}

			if ve.ValidationType == "error" {
				report.Errors = append(report.Errors, Issue{
					Message:  msg,
					Location: location,
					Rule:     "openapi-schema",
				})
			} else {
				report.Warnings = append(report.Warnings, Issue{
					Message:  msg,
					Location: location,
					Rule:     "openapi-schema",
				})
			}
		}
	}

	report.Valid = len(report.Errors) == 0

	return report
}

func writeReport(cfg Config, report *Report) error {
	if err := os.MkdirAll(cfg.OutputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	reportBytes, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	reportPath := filepath.Join(cfg.OutputDir, cfg.ReportFileName)
	if err := os.WriteFile(reportPath, reportBytes, 0o644); err != nil {
		return fmt.Errorf("failed to write report file '%s': %w", reportPath, err)
	}

	return nil
}
