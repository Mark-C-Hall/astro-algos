package angle

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

// Chapter 1's example of a large angle: sin 36000030° = sin 30°.
func TestTrigFunctionOfLargeAngle(t *testing.T) {
	angle := FromDegrees(36000030)
	got := math.Sin(angle.Radians())
	if !floatsEqual(got, 0.5, bookTolerance) {
		t.Errorf("Sin(angle) = %.10f; want 0.5", got)
	}
}

// Chapter 2's quick checks of a language's trig functions.
func TestTrigAccuracy(t *testing.T) {
	const tolerance = 1e-15
	tests := []struct {
		name string
		got  float64
		want float64
	}{
		{"sin 0.61 rad", FromRadians(0.61).Sin(), 0.5728674601004813},
		{"4 atan 1", 4 * Atan2(1, 1).Radians(), math.Pi},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !floatsEqual(tt.got, tt.want, tolerance) {
				t.Errorf("%s = %.16f; want %.16f", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestDMSDegrees(t *testing.T) {
	tests := []struct {
		name string
		in   DMS
		want float64
	}{
		{"obliquity", DMS{Deg: 23, Min: 26, Sec: 49}, 23.44694444},
		{"negative", DMS{Neg: true, Deg: 13, Min: 47, Sec: 22}, -13.78944444},
		{"negative zero degrees", DMS{Neg: true, Min: 32, Sec: 41}, -0.54472222},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.in.Degrees(); !floatsEqual(got, tt.want, bookTolerance) {
				t.Errorf("%+v.Degrees() = %.8f; want %.8f", tt.in, got, tt.want)
			}
		})
	}
}

func TestHMSDegrees(t *testing.T) {
	tests := []struct {
		name string
		in   HMS
		want float64
	}{
		{"Normal", HMS{Hr: 9, Min: 14, Sec: 55.8}, 138.73250},
		{"negative", HMS{Neg: true, Hr: 1, Min: 30}, -22.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.in.Degrees(); !floatsEqual(got, tt.want, bookTolerance) {
				t.Errorf("%+v.Degrees() = %.8f; want %.8f", tt.in, got, tt.want)
			}
		})
	}
}

func TestFromDegrees(t *testing.T) {
	tests := []struct {
		name string
		in   float64
		want float64
		tol  float64
	}{
		{"half turn", 180, math.Pi, bookTolerance},
		{"negative", -90, -math.Pi / 2, bookTolerance},
		{"reduced mod 360", 765, math.Pi / 4, bookTolerance},
		{"moon mean longitude reduced", 192725.46868893, 125.46868893 * math.Pi / 180, bookTolerance},
		{"moon mean longitude as printed", 192725.469, 125.469 * math.Pi / 180, printed(3) * math.Pi / 180},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromDegrees(tt.in).Radians(); !floatsEqual(got, tt.want, tt.tol) {
				t.Errorf("FromDegrees(%v).Radians() = %.8f; want %.8f", tt.in, got, tt.want)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name string
		in   Angle
		want Angle
	}{
		{"zero", 0, 0},
		{"in range", 1, 1},
		{"full turn", 2 * math.Pi, 0},
		{"negative", -math.Pi / 2, 3 * math.Pi / 2},
		{"multiple turns", 5 * math.Pi, math.Pi},
		{"tiny negative rounds to full turn", -1e-18, 0}, // -1e-18 + 2π rounds to 2π, which wraps to 0
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.in.Normalize(); !floatsEqual(float64(got), float64(tt.want), bookTolerance) {
				t.Errorf("Angle(%v).Normalize() = %.10f; want %.10f", tt.in, got, tt.want)
			}
		})
	}
}

func TestNewDMS(t *testing.T) {
	tests := []struct {
		name    string
		neg     bool
		d, m    int
		s       float64
		wantErr bool
	}{
		{"valid", false, 23, 26, 49, false},
		{"valid negative", true, 0, 59, 59.999, false},
		{"negative degrees", false, -1, 0, 0, true},
		{"negative minutes", false, 0, -1, 0, true},
		{"negative seconds", false, 0, 0, -1, true},
		{"minutes too large", false, 0, 60, 0, true},
		{"seconds too large", false, 0, 0, 60, true},
		{"NaN seconds", false, 0, 0, math.NaN(), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDMS(tt.neg, tt.d, tt.m, tt.s)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewDMS(%v, %d, %d, %v) error = %v; wantErr %v", tt.neg, tt.d, tt.m, tt.s, err, tt.wantErr)
			}
			if err != nil && got != (DMS{}) {
				t.Errorf("NewDMS returned %+v with error; want zero value", got)
			}
		})
	}
}

func TestNewHMS(t *testing.T) {
	tests := []struct {
		name    string
		neg     bool
		h, m    int
		s       float64
		wantErr bool
	}{
		{"valid", false, 9, 14, 55.8, false},
		{"negative hour angle", true, 30, 0, 0, false},
		{"negative field", false, 0, -1, 0, true},
		{"minutes too large", false, 0, 60, 0, true},
		{"NaN seconds", false, 0, 0, math.NaN(), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHMS(tt.neg, tt.h, tt.m, tt.s)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewHMS(%v, %d, %d, %v) error = %v; wantErr %v", tt.neg, tt.h, tt.m, tt.s, err, tt.wantErr)
			}
			if err != nil && got != (HMS{}) {
				t.Errorf("NewHMS returned %+v with error; want zero value", got)
			}
		})
	}
}

func TestHMSValidateRA(t *testing.T) {
	tests := []struct {
		name    string
		in      HMS
		wantErr bool
	}{
		{"valid", HMS{Hr: 23, Min: 59, Sec: 59.9}, false},
		{"negative", HMS{Neg: true, Hr: 1}, true},
		{"hours too large", HMS{Hr: 24}, true},
		{"invalid minutes", HMS{Hr: 1, Min: 60}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.in.ValidateRA(); (err != nil) != tt.wantErr {
				t.Errorf("%+v.ValidateRA() error = %v; wantErr %v", tt.in, err, tt.wantErr)
			}
		})
	}
}

func TestAsinAcos(t *testing.T) {
	tests := []struct {
		name     string
		in       float64
		wantAsin float64
		wantAcos float64
	}{
		{"zero", 0, 0, math.Pi / 2},
		{"half", 0.5, math.Pi / 6, math.Pi / 3},
		{"one", 1, math.Pi / 2, 0},
		{"minus one", -1, -math.Pi / 2, math.Pi},
		{"rounding above one", 1 + 1e-15, math.Pi / 2, 0},
		{"rounding below minus one", -1 - 1e-15, -math.Pi / 2, math.Pi},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Asin(tt.in).Radians(); !floatsEqual(got, tt.wantAsin, bookTolerance) {
				t.Errorf("Asin(%v) = %.10f; want %.10f", tt.in, got, tt.wantAsin)
			}
			if got := Acos(tt.in).Radians(); !floatsEqual(got, tt.wantAcos, bookTolerance) {
				t.Errorf("Acos(%v) = %.10f; want %.10f", tt.in, got, tt.wantAcos)
			}
		})
	}
}

func TestAsinAcosPanicOutOfRange(t *testing.T) {
	for _, x := range []float64{1 + 1e-11, -1 - 1e-11, 1.5, math.NaN()} {
		assertPanics(t, "Asin", func() { Asin(x) })
		assertPanics(t, "Acos", func() { Acos(x) })
	}
}

func assertPanics(t *testing.T, name string, f func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("%s did not panic", name)
		}
	}()
	f()
}
