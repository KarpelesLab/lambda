package lambda

import (
	"testing"
)

func TestISPRIME(t *testing.T) {
	tests := []struct {
		n    int
		want bool
	}{
		{2, true},
		{3, true},
		{4, false},
		{5, true},
		{6, false},
		{7, true},
		{8, false},
		{9, false},
		{10, false},
		{11, true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := Application{
				Func: IS_PRIME,
				Arg:  ChurchNumeral(tt.n),
			}

			// IS_PRIME is extremely computationally intensive
			reduced, steps := Reduce(result, 50000)
			got := ToBool(reduced)
			t.Logf("IS_PRIME %d = %v (want %v) in %d steps", tt.n, got, tt.want, steps)

			if got != tt.want {
				t.Errorf("IS_PRIME %d = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}