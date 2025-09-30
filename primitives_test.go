package lambda

import (
	"testing"
)

func TestEQ(t *testing.T) {
	tests := []struct {
		m, n int
		want bool
	}{
		{5, 5, true},
		{3, 5, false},
		{5, 3, false},
		{0, 0, true},
		{10, 10, true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := Application{
				Func: Application{
					Func: EQ,
					Arg:  ChurchNumeral(tt.m),
				},
				Arg: ChurchNumeral(tt.n),
			}

			reduced, _ := Reduce(result, 5000)
			got := ToBool(reduced)

			if got != tt.want {
				t.Errorf("EQ %d %d = %v, want %v", tt.m, tt.n, got, tt.want)
			}
		})
	}
}

func TestMAX(t *testing.T) {
	tests := []struct {
		a, b int
		want int
	}{
		{5, 3, 5},
		{3, 5, 5},
		{7, 7, 7},
		{0, 5, 5},
		{5, 0, 5},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := Application{
				Func: Application{
					Func: MAX,
					Arg:  ChurchNumeral(tt.a),
				},
				Arg: ChurchNumeral(tt.b),
			}

			reduced, _ := Reduce(result, 5000)
			got := ToInt(reduced)

			if got != tt.want {
				t.Errorf("MAX %d %d = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestMIN(t *testing.T) {
	tests := []struct {
		a, b int
		want int
	}{
		{5, 3, 3},
		{3, 5, 3},
		{7, 7, 7},
		{0, 5, 0},
		{5, 0, 0},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := Application{
				Func: Application{
					Func: MIN,
					Arg:  ChurchNumeral(tt.a),
				},
				Arg: ChurchNumeral(tt.b),
			}

			reduced, _ := Reduce(result, 5000)
			got := ToInt(reduced)

			if got != tt.want {
				t.Errorf("MIN %d %d = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestGCD(t *testing.T) {
	tests := []struct {
		a, b int
		want int
	}{
		{10, 5, 5},
		{6, 4, 2},
		{9, 6, 3},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := Application{
				Func: Application{
					Func: GCD,
					Arg:  ChurchNumeral(tt.a),
				},
				Arg: ChurchNumeral(tt.b),
			}

			reduced, steps := Reduce(result, 20000)
			got := ToInt(reduced)
			t.Logf("GCD %d %d = %d (want %d) in %d steps", tt.a, tt.b, got, tt.want, steps)

			if got != tt.want {
				t.Errorf("GCD %d %d = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}