package lambda

import (
	"os"
	"testing"
)

func TestGenerateDiagrams(t *testing.T) {
	// Create diagrams directory if it doesn't exist
	if err := os.MkdirAll("diagrams", 0755); err != nil {
		t.Fatal(err)
	}

	// Generate diagrams for various lambda terms
	terms := map[string]Object{
		"identity":    I,
		"k":           K,
		"s":           S,
		"b":           B,
		"c":           C,
		"w":           W,
		"church_0":    ChurchNumeral(0),
		"church_1":    ChurchNumeral(1),
		"church_2":    ChurchNumeral(2),
		"church_3":    ChurchNumeral(3),
		"true":        TRUE,
		"false":       FALSE,
		"omega":       U,
		"application": Application{Func: I, Arg: Var{Name: "x"}},
	}

	for name, term := range terms {
		svg := ToDiagram(term).ToSVG()

		// Read existing file if it exists
		filename := "diagrams/" + name + ".svg"
		existing, err := os.ReadFile(filename)

		// Write the file
		if err := os.WriteFile(filename, []byte(svg), 0644); err != nil {
			t.Fatal(err)
		}

		// If file existed and content changed, report it
		if err == nil && string(existing) != svg {
			t.Logf("Updated diagram: %s", filename)
		} else if err != nil {
			t.Logf("Created diagram: %s", filename)
		}
	}
}