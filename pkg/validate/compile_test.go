package validate_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/kalo-build/go-util/assertfile"
	"github.com/kalo-build/plugin-openapi-validator/pkg/validate"
	"github.com/stretchr/testify/suite"
)

type CompileTestSuite struct {
	assertfile.FileSuite

	TestDataPath    string
	GroundTruthPath string
}

func TestCompileTestSuite(t *testing.T) {
	suite.Run(t, new(CompileTestSuite))
}

func (suite *CompileTestSuite) SetupTest() {
	_, filename, _, _ := runtime.Caller(0)
	suite.TestDataPath = filepath.Join(filepath.Dir(filename), "..", "..", "testdata")
	suite.GroundTruthPath = filepath.Join(suite.TestDataPath, "ground-truth")
}

func (suite *CompileTestSuite) TestValidateSpec_Valid() {
	inputDir := filepath.Join(suite.TestDataPath, "input")
	outputDir := suite.T().TempDir()

	cfg := validate.Config{
		InputDir:       inputDir,
		OutputDir:      outputDir,
		SpecFileName:   "valid.yaml",
		ReportFileName: "valid-report.json",
	}

	report, err := validate.ValidateSpec(cfg)
	suite.NoError(err)
	suite.True(report.Valid)
	suite.Empty(report.Errors)

	reportPath := filepath.Join(outputDir, "valid-report.json")
	gtReportPath := filepath.Join(suite.GroundTruthPath, "valid-report.json")
	suite.FileExists(reportPath)
	suite.FileEquals(reportPath, gtReportPath)
}
