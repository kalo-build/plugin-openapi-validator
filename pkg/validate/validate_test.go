package validate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testdataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "testdata")
}

func TestValidateSpec_ValidSpec(t *testing.T) {
	inputDir := filepath.Join(testdataDir(), "input")
	outputDir := t.TempDir()

	cfg := Config{
		InputDir:       inputDir,
		OutputDir:      outputDir,
		SpecFileName:   "valid.yaml",
		ReportFileName: "report.json",
	}

	report, err := ValidateSpec(cfg)
	require.NoError(t, err)
	assert.True(t, report.Valid)
	assert.Empty(t, report.Errors)

	reportPath := filepath.Join(outputDir, "report.json")
	reportBytes, err := os.ReadFile(reportPath)
	require.NoError(t, err)

	var diskReport Report
	err = json.Unmarshal(reportBytes, &diskReport)
	require.NoError(t, err)
	assert.True(t, diskReport.Valid)
}

func TestValidateBytes_ValidDocument(t *testing.T) {
	specBytes := []byte(`openapi: "3.1.0"
info:
  title: Test
  version: "1.0.0"
paths:
  /test:
    get:
      operationId: getTest
      responses:
        "200":
          description: OK
`)
	report := ValidateBytes(specBytes, "test.yaml")
	assert.True(t, report.Valid)
	assert.Empty(t, report.Errors)
	assert.Equal(t, "test.yaml", report.SpecFile)
}

func TestValidateBytes_InvalidYAML(t *testing.T) {
	specBytes := []byte(`{not valid yaml: [}`)
	report := ValidateBytes(specBytes, "bad.yaml")
	assert.False(t, report.Valid)
	assert.NotEmpty(t, report.Errors)
}

func TestConfig_Resolve(t *testing.T) {
	cfg := Config{}
	cfg.Resolve()
	assert.Equal(t, "openapi.yaml", cfg.SpecFileName)
	assert.Equal(t, "validation-report.json", cfg.ReportFileName)
}
