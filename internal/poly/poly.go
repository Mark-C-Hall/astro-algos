// Package poly provides polynomial evaluation.
package poly

import "slices"

// Horner evaluates coeffs[0] + coeffs[1]*x + coeffs[2]*x² + ... using Horner's
// method. Coefficients are in ascending order of power.
func Horner(x float64, coeffs ...float64) float64 {
	var y float64
	for _, coeff := range slices.Backward(coeffs) {
		y = y*x + coeff
	}
	return y
}
