package advmath

import "math"

func Smoothstep(edge0 float64, edge1 float64, x float64) float64 {
	// Scale, bias and saturate x to 0..1 range
	x = ClampFloat64((x-edge0)/(edge1-edge0), 0.0, 1.0)
	// Evaluate polynomial
	return x * x * (3 - 2*x)
}

func InverseLerpFloat64(min float64, max float64, x float64) float64 {
	return (x - min) / (max - min)
}

func ClampFloat64(d float64, min float64, max float64) float64 {
	var t float64
	if d < min {
		t = min
	} else {
		t = d
	}
	if t > max {
		return max
	}
	return t
}

func CircIn(completed float64) float64 {
	return 1 - math.Sqrt(1-completed*completed)
}
