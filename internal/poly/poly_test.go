package poly

import (
	"math"
	"testing"
)

const bookTolerance = 1e-6

func floatsEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

func TestHorner(t *testing.T) {
	tests := []struct {
		name   string
		x      float64
		coeffs []float64
		want   float64
	}{
		{"no coefficients", 2, nil, 0},
		{"constant only", 7, []float64{3}, 3},
		{"worked example", 2, []float64{1, 2, 3, 4, 5}, 129},
		{"x is zero", 0, []float64{1, 2, 3, 4, 5}, 1},
		{"negative x", -1, []float64{1, 2, 3}, 2},
		{"uranus at J2000", 0, []float64{0.04638122, -0.000027293, 0.000000078972}, 0.04638122},
		{"uranus far past", -3307.9, []float64{0.04638122, -0.000027293, 0.000000078972}, 1.00079134},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Horner(tt.x, tt.coeffs...); !floatsEqual(got, tt.want, bookTolerance) {
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
