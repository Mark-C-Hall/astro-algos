package poly

import (
	"math"
	"testing"
)

// bookTolerance suits expected values written to 8 decimals.
const bookTolerance = 1e-6

func floatsEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

// printed returns half a unit in the last place of a value printed to n decimals.
func printed(n int) float64 { return 0.5 * math.Pow10(-n) }

func TestHorner(t *testing.T) {
	tests := []struct {
		name   string
		x      float64
		coeffs []float64
		want   float64
		tol    float64
	}{
		{"no coefficients", 2, nil, 0, bookTolerance},
		{"constant only", 7, []float64{3}, 3, bookTolerance},
		{"worked example", 2, []float64{1, 2, 3, 4, 5}, 129, bookTolerance},
		{"x is zero", 0, []float64{1, 2, 3, 4, 5}, 1, bookTolerance},
		{"negative x", -1, []float64{1, 2, 3}, 2, bookTolerance},
		// Uranus's orbital eccentricity, which exceeds 1 far outside its valid range.
		{"uranus at J2000", 0, []float64{0.04638122, -0.000027293, 0.000000078972}, 0.04638122, bookTolerance},
		{"uranus far past", -3307.9, []float64{0.04638122, -0.000027293, 0.000000078972}, 1.00079134, bookTolerance},
		{"moon mean longitude T=0.4", 0.4, []float64{218.3164477, 481267.88123421, -0.0015786, 0.0000019}, 192725.46868893, bookTolerance},
		{"moon mean longitude as printed", 0.4, []float64{218.3164477, 481267.88123421, -0.0015786, 0.0000019}, 192725.469, printed(3)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Horner(tt.x, tt.coeffs...); !floatsEqual(got, tt.want, tt.tol) {
				t.Errorf("Horner(%v, %v...) = %.8f; want %.8f", tt.x, tt.coeffs, got, tt.want)
			}
		})
	}
}

// naive evaluates the polynomial term by term with math.Pow,
// the approach Meeus contrasts with Horner's method in Chapter 1.
func naive(x float64, coeffs ...float64) float64 {
	var y float64
	for i, c := range coeffs {
		y += c * math.Pow(x, float64(i))
	}
	return y
}

// Package-level inputs keep the compiler from constant-folding the calls,
// and sink keeps the results "used" so the work isn't optimized away.
var (
	benchX      = 1.23
	benchCoeffs = []float64{1, 2, 3, 4, 5}
	sink        float64
)

func BenchmarkHorner(b *testing.B) {
	for b.Loop() {
		sink = Horner(benchX, benchCoeffs...)
	}
}

func BenchmarkNaive(b *testing.B) {
	for b.Loop() {
		sink = naive(benchX, benchCoeffs...)
	}
}
