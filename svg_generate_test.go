package lambda

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGenerateExampleSVGs(t *testing.T) {
	// Find repo root relative to this test file
	_, thisFile, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(thisFile), "examples")

	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}

	opts := &SVGOptions{
		CellSize:   20,
		Padding:    10,
		Background: "#1a1a2e",
		Saturation: 0.8,
		Value:      0.95,
	}

	terms := []struct {
		name string
		term Term
	}{
		{"identity", I},
		{"K", K},
		{"S", S},
		{"omega", OMEGA},
		{"Y", Y},
		{"church2", ChurchNumeral(2)},
		{"church3", ChurchNumeral(3)},
	}

	for _, tt := range terms {
		svg := DiagramSVG(tt.term, opts)
		path := filepath.Join(dir, tt.name+".svg")
		if err := os.WriteFile(path, []byte(svg), 0644); err != nil {
			t.Errorf("failed to write %s: %v", path, err)
			continue
		}
		t.Logf("Wrote %s", path)
	}
}
