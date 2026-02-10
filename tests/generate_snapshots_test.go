package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/brownhounds/openapi-tsgen/schema"
)

func TestGenerateSnapshotsMatchFixtures(t *testing.T) {
	fixturesDir := "fixtures"
	snapshotsDir := "snapshots"
	tmpDir := ".generated"

	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		t.Fatalf("create tmp dir: %v", err)
	}

	oldNow := schema.Now
	oldVersion := schema.CLIVersion
	schema.Now = func() time.Time {
		return time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC)
	}
	schema.CLIVersion = "dev"
	defer func() {
		schema.Now = oldNow
		schema.CLIVersion = oldVersion
	}()

	entries, err := os.ReadDir(fixturesDir)
	if err != nil {
		t.Fatalf("read fixtures dir: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".fixture.yml") && !strings.HasSuffix(name, ".fixture.json") {
			continue
		}

		fixturePath := filepath.Join(fixturesDir, name)
		base := strings.TrimSuffix(strings.TrimSuffix(name, ".fixture.yml"), ".fixture.json")
		format := schema.InputYAML
		snapshotName := base + ".yml.snapshot.ts"
		if strings.HasSuffix(name, ".fixture.json") {
			format = schema.InputJSON
			snapshotName = base + ".json.snapshot.ts"
		}

		outPath := filepath.Join(tmpDir, snapshotName)
		if err := schema.WriteSchema(fixturePath, outPath, format); err != nil {
			t.Fatalf("generate %s: %v", name, err)
		}

		expectedPath := filepath.Join(snapshotsDir, snapshotName)
		expected, err := os.ReadFile(expectedPath)
		if err != nil {
			t.Fatalf("read %s: %v", expectedPath, err)
		}
		got, err := os.ReadFile(outPath)
		if err != nil {
			t.Fatalf("read %s: %v", outPath, err)
		}
		if normalizeSnapshot(string(expected)) != normalizeSnapshot(string(got)) {
			t.Fatalf("snapshot mismatch: %s vs %s\n%s", snapshotName, outPath, diffText(normalizeSnapshot(string(expected)), normalizeSnapshot(string(got))))
		}
	}
}
