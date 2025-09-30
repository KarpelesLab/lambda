package lambda

import (
	"testing"
)

func TestDIV2(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{0, 0},
		{1, 0},
		{2, 1},
		{3, 1},
		{4, 2},
		{5, 2},
		{10, 5},
		{11, 5},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := Application{
				Func: DIV2,
				Arg:  ChurchNumeral(tt.n),
			}

			reduced, _ := Reduce(result, 5000)
			got := ToInt(reduced)

			if got != tt.want {
				t.Errorf("DIV2 %d = %d, want %d", tt.n, got, tt.want)
			}
		})
	}
}

func TestISODD(t *testing.T) {
	tests := []struct {
		n    int
		want bool
	}{
		{0, false},
		{1, true},
		{2, false},
		{3, true},
		{4, false},
		{5, true},
		{10, false},
		{11, true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := Application{
				Func: ISODD,
				Arg:  ChurchNumeral(tt.n),
			}

			reduced, _ := Reduce(result, 5000)
			got := ToBool(reduced)

			if got != tt.want {
				t.Errorf("ISODD %d = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

func TestISEVEN(t *testing.T) {
	tests := []struct {
		n    int
		want bool
	}{
		{0, true},
		{1, false},
		{2, true},
		{3, false},
		{4, true},
		{5, false},
		{10, true},
		{11, false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := Application{
				Func: ISEVEN,
				Arg:  ChurchNumeral(tt.n),
			}

			reduced, _ := Reduce(result, 5000)
			got := ToBool(reduced)

			if got != tt.want {
				t.Errorf("ISEVEN %d = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

func TestPOWMOD(t *testing.T) {
	tests := []struct {
		a, e, m int
		want    int
	}{
		{2, 0, 5, 1},     // 2^0 mod 5 = 1
		{2, 1, 5, 2},     // 2^1 mod 5 = 2
		{2, 2, 5, 4},     // 2^2 mod 5 = 4
		{3, 2, 7, 2},     // 3^2 mod 7 = 9 mod 7 = 2
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := Application{
				Func: Application{
					Func: Application{
						Func: POWMOD,
						Arg:  ChurchNumeral(tt.a),
					},
					Arg: ChurchNumeral(tt.e),
				},
				Arg: ChurchNumeral(tt.m),
			}

			reduced, steps := Reduce(result, 20000)
			got := ToInt(reduced)
			t.Logf("POWMOD %d %d %d = %d (want %d) in %d steps", tt.a, tt.e, tt.m, got, tt.want, steps)

			if got != tt.want {
				t.Errorf("POWMOD %d %d %d = %d, want %d", tt.a, tt.e, tt.m, got, tt.want)
			}
		})
	}
}

func TestPOWMOD_PRIME(t *testing.T) {
	tests := []struct {
		a, e, m int
		want    int
	}{
		{2, 0, 5, 1},     // 2^0 mod 5 = 1
		{2, 1, 5, 2},     // 2^1 mod 5 = 2
		{2, 2, 5, 4},     // 2^2 mod 5 = 4
		{3, 2, 7, 2},     // 3^2 mod 7 = 9 mod 7 = 2
		{2, 3, 5, 3},     // 2^3 mod 5 = 8 mod 5 = 3
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			// POWMOD_PRIME takes 4 args: a, e, m, r (where r is the accumulator, initially 1)
			result := Application{
				Func: Application{
					Func: Application{
						Func: Application{
							Func: POWMOD_PRIME,
							Arg:  ChurchNumeral(tt.a),
						},
						Arg: ChurchNumeral(tt.e),
					},
					Arg: ChurchNumeral(tt.m),
				},
				Arg: ONE, // Initial accumulator
			}

			reduced, steps := Reduce(result, 20000)
			got := ToInt(reduced)
			t.Logf("POWMOD_PRIME %d %d %d = %d (want %d) in %d steps", tt.a, tt.e, tt.m, got, tt.want, steps)

			if got != tt.want {
				t.Errorf("POWMOD_PRIME %d %d %d = %d, want %d", tt.a, tt.e, tt.m, got, tt.want)
			}
		})
	}
}