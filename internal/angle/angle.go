// Package angle provides angle types and conversions for astronomical calculations.
package angle

import (
	"errors"
	"fmt"
	"math"
)

// Angle is an angle in radians.
type Angle float64

// FromDegrees converts d to an Angle, first reducing it modulo 360 (keeping
// its sign) so very large inputs don't lose precision.
func FromDegrees(d float64) Angle { return Angle(math.Mod(d, 360) * math.Pi / 180) }

func FromRadians(r float64) Angle { return Angle(r) }
func (a Angle) Degrees() float64  { return float64(a) * 180 / math.Pi }
func (a Angle) Radians() float64  { return float64(a) }

func (a Angle) Sin() float64 { return math.Sin(float64(a)) }
func (a Angle) Cos() float64 { return math.Cos(float64(a)) }
func (a Angle) Tan() float64 { return math.Tan(float64(a)) }

// Asin and Acos clamp inputs that drift slightly past ±1 from rounding error,
// and panic on anything further out (or NaN), which indicates a bug upstream.
func Asin(x float64) Angle     { return Angle(math.Asin(clampUnit(x))) }
func Acos(x float64) Angle     { return Angle(math.Acos(clampUnit(x))) }
func Atan2(y, x float64) Angle { return Angle(math.Atan2(y, x)) }

// clampTolerance is how far past ±1 an inverse trig input may drift from
// rounding error before it is treated as a bug.
const clampTolerance = 1e-12

func clampUnit(x float64) float64 {
	if math.IsNaN(x) || math.Abs(x) > 1+clampTolerance {
		panic(fmt.Sprintf("angle: inverse trig input %v out of range [-1, 1]", x))
	}
	return math.Max(-1, math.Min(1, x))
}

// Normalize reduces a to the range [0, 2π).
func (a Angle) Normalize() Angle {
	r := math.Mod(float64(a), 2*math.Pi)
	if r < 0 {
		r += 2 * math.Pi
	}
	if r >= 2*math.Pi {
		r = 0
	}
	return Angle(r)
}

// DMS is an angle in degrees, minutes and seconds of arc. Fields are
// non-negative and the sign is held in Neg, so -0°32′41″ is representable.
type DMS struct {
	Neg bool
	Deg int
	Min int
	Sec float64
}

// Validate checks that all fields are non-negative and Min and Sec are below 60.
func (d DMS) Validate() error {
	switch {
	case d.Deg < 0 || d.Min < 0 || d.Sec < 0:
		return errors.New("angle: fields must be non-negative; use Neg for sign")
	case d.Min >= 60:
		return fmt.Errorf("angle: minutes %d out of range [0, 60)", d.Min)
	case math.IsNaN(d.Sec) || d.Sec >= 60:
		return fmt.Errorf("angle: seconds %v out of range [0, 60)", d.Sec)
	}
	return nil
}

// NewDMS returns a validated DMS.
func NewDMS(neg bool, d, m int, s float64) (DMS, error) {
	dms := DMS{
		Neg: neg,
		Deg: d,
		Min: m,
		Sec: s,
	}
	if err := dms.Validate(); err != nil {
		return DMS{}, err
	}
	return dms, nil
}

func (d DMS) Degrees() float64 {
	result := float64(d.Deg) + float64(d.Min)/60 + d.Sec/3600
	if d.Neg {
		return -result
	}
	return result
}

func (d DMS) Angle() Angle {
	return FromDegrees(d.Degrees())
}

// HMS is an angle in hours, minutes and seconds of time, where 1h = 15°.
// As with DMS, the sign is held in Neg (e.g. for hour angles).
type HMS struct {
	Neg bool
	Hr  int
	Min int
	Sec float64
}

// Validate checks that all fields are non-negative and Min and Sec are below 60.
func (h HMS) Validate() error {
	switch {
	case h.Hr < 0 || h.Min < 0 || h.Sec < 0:
		return errors.New("angle: fields must be non-negative; use Neg for sign")
	case h.Min >= 60:
		return fmt.Errorf("angle: minutes %d out of range [0, 60)", h.Min)
	case math.IsNaN(h.Sec) || h.Sec >= 60:
		return fmt.Errorf("angle: seconds %v out of range [0, 60)", h.Sec)
	}
	return nil
}

// ValidateRA additionally requires a right ascension: non-negative with Hr below 24.
func (h HMS) ValidateRA() error {
	if err := h.Validate(); err != nil {
		return err
	}
	if h.Neg {
		return errors.New("angle: right ascension cannot be negative")
	}
	if h.Hr >= 24 {
		return fmt.Errorf("angle: hours %d out of range [0, 24)", h.Hr)
	}
	return nil
}

// NewHMS returns a validated HMS. Use NewRA for right ascension.
func NewHMS(neg bool, h, m int, s float64) (HMS, error) {
	hms := HMS{
		Neg: neg,
		Hr:  h,
		Min: m,
		Sec: s,
	}
	if err := hms.Validate(); err != nil {
		return HMS{}, err
	}
	return hms, nil
}

// NewRA returns an HMS validated as a right ascension.
func NewRA(h, m int, s float64) (HMS, error) {
	hms := HMS{
		Neg: false,
		Hr:  h,
		Min: m,
		Sec: s,
	}
	if err := hms.ValidateRA(); err != nil {
		return HMS{}, err
	}
	return hms, nil
}

func (h HMS) Degrees() float64 {
	result := 15 * (float64(h.Hr) + float64(h.Min)/60 + h.Sec/3600)
	if h.Neg {
		return -result
	}
	return result
}

func (h HMS) Angle() Angle {
	return FromDegrees(h.Degrees())
}
