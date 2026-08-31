# plugin-openapi-validator

Kalo plugin that validates an OpenAPI 3.1 specification using [libopenapi](https://github.com/pb33f/libopenapi) and [libopenapi-validator](https://github.com/pb33f/libopenapi-validator). Produces a validation report and fails the pipeline on errors.

## Overview

This plugin reads a `KA:OA1:YAML1` OpenAPI specification, validates it against the OpenAPI 3.1 standard, and writes a JSON validation report. If the spec contains errors, the plugin exits with a non-zero code, halting the pipeline.

Use this as a quality gate after OpenAPI generation to catch malformed specs before they reach downstream consumers (docs, SDK generation, etc.).

## Input / Output

| Direction | Format | Description |
|-----------|--------|-------------|
| **Input** | `KA:OA1:YAML1` | OpenAPI 3.1 YAML specification |
| **Output** | `KA:OA1:REPORT1` | Validation report (JSON) |

## Configuration

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `specFileName` | string | `"openapi.yaml"` | Name of the spec file to read from input |
| `failOnWarning` | boolean | `false` | Fail the pipeline on warnings (not just errors) |
| `reportFileName` | string | `"validation-report.json"` | Output report file name |

## Usage in `kalo.yaml`

```yaml
stores:
  KA_OA_YAML:
    type: localFileSystem
    path: ./docs/openapi
  KA_OA_REPORT:
    type: localFileSystem
    path: ./docs/openapi-report

config:
  "@kalo-build/plugin-openapi-validator":
    failOnWarning: true

pipeline:
  compile:
    # ... generate openapi.yaml first ...
    openapi-validate:
      plugin: "@kalo-build/plugin-openapi-validator"

plugins:
  "@kalo-build/plugin-openapi-validator":
    input:
      store: KA_OA_YAML
    output:
      format: "KA:OA1:REPORT1"
      store: KA_OA_REPORT
```

## Validation Report Format

```json
{
  "valid": true,
  "errors": [],
  "warnings": [
    {
      "message": "Description of warning",
      "location": "path/to/issue",
      "rule": "openapi-schema"
    }
  ],
  "specFile": "openapi.yaml",
  "openapiVersion": "1.0.0"
}
```

## How It Works

1. Reads the OpenAPI YAML spec from the input directory
2. Parses the document using `libopenapi.NewDocument()`
3. Builds the V3 model to detect structural issues
4. Runs `validator.ValidateDocument()` for full schema validation
5. Writes a JSON report to the output directory
6. Exits with code `0` (valid) or `2` (errors found)

## Wrapped Packages

| Package | Purpose |
|---------|---------|
| [pb33f/libopenapi](https://github.com/pb33f/libopenapi) | OpenAPI 3.1 document parsing and model building |
| [pb33f/libopenapi-validator](https://github.com/pb33f/libopenapi-validator) | Full OpenAPI document validation |

## Development

```bash
go test -v ./...
```

## Project Structure

```
plugin-openapi-validator/
├── cmd/plugin/main.go           # WASM entrypoint
├── pkg/validate/
│   ├── config.go                # Configuration with defaults
│   ├── report.go                # Report types
│   ├── validate.go              # Validation logic (libopenapi wrapper)
│   └── validate_test.go         # Tests
├── testdata/input/              # Test fixtures (valid + invalid specs)
├── plugin.yaml                  # Plugin manifest
└── README.md
```

## License

MIT
