package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kalo-build/plugin-openapi-validator/pkg/validate"
)

type StoreConfig struct {
	ID        uint32 `json:"id"`
	Type      string `json:"type"`
	MountPath string `json:"mountPath,omitempty"`
}

type PluginConfig struct {
	Stores     map[string]StoreConfig `json:"stores,omitempty"`
	InputPath  string                 `json:"inputPath,omitempty"`
	OutputPath string                 `json:"outputPath,omitempty"`
	Config     PluginConfigFields     `json:"config"`
	Verbose    bool                   `json:"verbose,omitempty"`
}

type PluginConfigFields struct {
	SpecFileName   string `json:"specFileName,omitempty"`
	FailOnWarning  bool   `json:"failOnWarning,omitempty"`
	ReportFileName string `json:"reportFileName,omitempty"`
}

const (
	ErrMissingConfig      = 3
	ErrInvalidConfig      = 4
	ErrInputPathRequired  = 12
	ErrOutputPathRequired = 13
	ErrValidationFailed   = 1
	ErrValidationErrors   = 2
)

func logInfo(verbose bool, format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stdout, format+"\n", args...)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: plugin-openapi-validator <config>")
		os.Exit(ErrMissingConfig)
	}

	rawConfig := os.Args[1]
	var pluginConfig PluginConfig
	if err := json.Unmarshal([]byte(rawConfig), &pluginConfig); err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing config JSON:", err)
		os.Exit(ErrInvalidConfig)
	}

	var inputPath, outputPath string

	if pluginConfig.Stores != nil {
		for _, store := range pluginConfig.Stores {
			switch store.MountPath {
			case "/input":
				inputPath = "/input"
			case "/output":
				outputPath = "/output"
			}
		}
	}

	if inputPath == "" && pluginConfig.InputPath != "" {
		inputPath = pluginConfig.InputPath
	}
	if outputPath == "" && pluginConfig.OutputPath != "" {
		outputPath = pluginConfig.OutputPath
	}

	if inputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: Input path is required (directory containing OpenAPI spec)")
		os.Exit(ErrInputPathRequired)
	}
	if outputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: Output path is required")
		os.Exit(ErrOutputPathRequired)
	}

	if inputPath[0] != '/' {
		if abs, err := filepath.Abs(inputPath); err == nil {
			inputPath = abs
		}
	}
	if outputPath[0] != '/' {
		if abs, err := filepath.Abs(outputPath); err == nil {
			outputPath = abs
		}
	}

	cfg := validate.Config{
		InputDir:       inputPath,
		OutputDir:      outputPath,
		SpecFileName:   pluginConfig.Config.SpecFileName,
		FailOnWarning:  pluginConfig.Config.FailOnWarning,
		ReportFileName: pluginConfig.Config.ReportFileName,
	}
	cfg.Resolve()

	logInfo(pluginConfig.Verbose, "Validating OpenAPI spec: '%s'", filepath.Join(cfg.InputDir, cfg.SpecFileName))

	report, err := validate.ValidateSpec(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Validation failed:", err)
		os.Exit(ErrValidationFailed)
	}

	logInfo(pluginConfig.Verbose, "Validation complete: %d errors, %d warnings",
		len(report.Errors), len(report.Warnings))

	if len(report.Errors) > 0 {
		fmt.Fprintf(os.Stderr, "Validation found %d error(s)\n", len(report.Errors))
		for _, e := range report.Errors {
			fmt.Fprintf(os.Stderr, "  ERROR: %s\n", e.Message)
		}
		os.Exit(ErrValidationErrors)
	}

	if cfg.FailOnWarning && len(report.Warnings) > 0 {
		fmt.Fprintf(os.Stderr, "Validation found %d warning(s) (failOnWarning=true)\n", len(report.Warnings))
		for _, w := range report.Warnings {
			fmt.Fprintf(os.Stderr, "  WARN: %s\n", w.Message)
		}
		os.Exit(ErrValidationErrors)
	}

	logInfo(pluginConfig.Verbose, "Validation passed")
	os.Exit(0)
}
