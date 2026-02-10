package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.yaml.in/yaml/v3"
)

var Now = time.Now

var CLIVersion string

var (
	ErrSchemaPathRequired = errors.New("schema path is required")
	ErrOutputPathRequired = errors.New("output path is required")
)

type InputFormat string

const (
	InputYAML InputFormat = "yaml"
	InputJSON InputFormat = "json"
)

func WriteSchema(schemaPath, outPath string, format InputFormat) error {
	if schemaPath == "" {
		return ErrSchemaPathRequired
	}
	if outPath == "" {
		return ErrOutputPathRequired
	}

	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("read schema %q: %w", schemaPath, err)
	}

	var doc Document
	switch format {
	case InputJSON:
		if err := json.Unmarshal(data, &doc); err != nil {
			return fmt.Errorf("unmarshal schema %q: %w", schemaPath, err)
		}
	default:
		if err := yaml.Unmarshal(data, &doc); err != nil {
			return fmt.Errorf("unmarshal schema %q: %w", schemaPath, err)
		}
	}

	ir, err := ToIR(&doc)
	if err != nil {
		return fmt.Errorf("build IR: %w", err)
	}

	out := normalizeGeneratedOutput(EmitTypesFromIRAt(ir, Now(), CLIVersion, doc.OpenAPI))

	if existing, err := os.ReadFile(outPath); err == nil {
		existingNormalized := normalizeGeneratedOutput(string(existing))
		if stripGeneratedHeader(existingNormalized) == stripGeneratedHeader(out) {
			return nil
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("read output %q: %w", outPath, err)
	}

	if dir := filepath.Dir(outPath); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create output dir %q: %w", dir, err)
		}
	}

	if err := os.WriteFile(outPath, []byte(out), 0o644); err != nil {
		return fmt.Errorf("write output %q: %w", outPath, err)
	}

	return nil
}

func stripGeneratedHeader(s string) string {
	if !strings.HasPrefix(s, headerStart()) {
		return s
	}
	end := strings.Index(s, "*/\n\n")
	if end == -1 {
		return s
	}
	return s[end+4:]
}

func normalizeGeneratedOutput(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	out := strings.Join(lines, "\n")
	out = strings.TrimRight(out, "\n")
	return out + "\n"
}
