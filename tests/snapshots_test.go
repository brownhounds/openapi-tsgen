package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSnapshotsMatchJSONAndYAML(t *testing.T) {
	dir := "snapshots"
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read snapshots dir: %v", err)
	}

	ymlFiles := map[string]struct{}{}
	jsonFiles := map[string]struct{}{}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		switch {
		case strings.HasSuffix(name, ".yml.snapshot.ts"):
			ymlFiles[name] = struct{}{}
		case strings.HasSuffix(name, ".json.snapshot.ts"):
			jsonFiles[name] = struct{}{}
		}
	}

	for ymlName := range ymlFiles {
		jsonName := strings.Replace(ymlName, ".yml.snapshot.ts", ".json.snapshot.ts", 1)
		if _, ok := jsonFiles[jsonName]; !ok {
			t.Fatalf("missing json snapshot for %s", ymlName)
		}

		ymlPath := filepath.Join(dir, ymlName)
		jsonPath := filepath.Join(dir, jsonName)

		ymlData, err := os.ReadFile(ymlPath)
		if err != nil {
			t.Fatalf("read %s: %v", ymlName, err)
		}
		jsonData, err := os.ReadFile(jsonPath)
		if err != nil {
			t.Fatalf("read %s: %v", jsonName, err)
		}
		ymlText := normalizeSnapshot(string(ymlData))
		jsonText := normalizeSnapshot(string(jsonData))
		if ymlText != jsonText {
			t.Fatalf("snapshot mismatch: %s vs %s\n%s", ymlName, jsonName, diffText(ymlText, jsonText))
		}
	}

	for jsonName := range jsonFiles {
		ymlName := strings.Replace(jsonName, ".json.snapshot.ts", ".yml.snapshot.ts", 1)
		if _, ok := ymlFiles[ymlName]; !ok {
			t.Fatalf("missing yml snapshot for %s", jsonName)
		}
	}
}

func normalizeSnapshot(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		line = strings.TrimRight(line, " \t")
		if strings.HasPrefix(line, " * Generated at: ") {
			line = " * Generated at: <normalized>"
		} else if strings.HasPrefix(line, " * OpenAPI version:") {
			line = " * OpenAPI version:"
		}
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}

func diffText(expected, actual string) string {
	expLines := strings.Split(expected, "\n")
	actLines := strings.Split(actual, "\n")

	edits := diffEdits(expLines, actLines)
	var b strings.Builder
	b.WriteString("\n\ndiff:\n")
	for _, e := range edits {
		switch e.kind {
		case diffDelete:
			b.WriteString(fmt.Sprintf("\x1b[31m-%4d %s\x1b[0m\n", e.aLine, e.text))
		case diffInsert:
			b.WriteString(fmt.Sprintf("\x1b[32m+%4d %s\x1b[0m\n", e.bLine, e.text))
		}
	}
	b.WriteString("\n\n")
	return b.String()
}

type diffKind int

const (
	diffEqual diffKind = iota
	diffDelete
	diffInsert
)

type diffEdit struct {
	text  string
	kind  diffKind
	aLine int
	bLine int
}

func diffEdits(a, b []string) []diffEdit {
	n := len(a)
	m := len(b)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			switch {
			case a[i] == b[j]:
				dp[i][j] = dp[i+1][j+1] + 1
			case dp[i+1][j] >= dp[i][j+1]:
				dp[i][j] = dp[i+1][j]
			default:
				dp[i][j] = dp[i][j+1]
			}
		}
	}

	edits := []diffEdit{}
	i, j := 0, 0
	for i < n && j < m {
		switch {
		case a[i] == b[j]:
			i++
			j++
		case dp[i+1][j] >= dp[i][j+1]:
			edits = append(edits, diffEdit{kind: diffDelete, aLine: i + 1, bLine: j + 1, text: a[i]})
			i++
		default:
			edits = append(edits, diffEdit{kind: diffInsert, aLine: i + 1, bLine: j + 1, text: b[j]})
			j++
		}
	}
	for i < n {
		edits = append(edits, diffEdit{kind: diffDelete, aLine: i + 1, bLine: j + 1, text: a[i]})
		i++
	}
	for j < m {
		edits = append(edits, diffEdit{kind: diffInsert, aLine: i + 1, bLine: j + 1, text: b[j]})
		j++
	}
	return edits
}
