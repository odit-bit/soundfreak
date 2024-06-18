package fx

import "math"

type FilterType int

const (
	LowPass FilterType = iota
	HighPass
	LowShelf
	HighShelf
)

// IIRFilter represents a second-order IIR filter
type IIRFilter struct {
	b0, b1, b2 float64
	a1, a2     float64
	x1, x2     float64
	y1, y2     float64
}

// NewIIRFilter creates a new IIR filter with the given coefficients
func NewIIRFilter(b0, b1, b2, a1, a2 float64) *IIRFilter {
	return &IIRFilter{
		b0: b0, b1: b1, b2: b2,
		a1: a1, a2: a2,
	}
}

// Process processes the input signal and returns the filtered output
func (f *IIRFilter) Process(input float64) float64 {
	output := f.b0*input + f.b1*f.x1 + f.b2*f.x2 - f.a1*f.y1 - f.a2*f.y2
	f.x2 = f.x1
	f.x1 = input
	f.y2 = f.y1
	f.y1 = output
	return output
}

// ComputeCoefficients calculates the coefficients for the given filter type and parameters
func ComputeCoefficients(filterType FilterType, cutoff, sampleRate, gainDB float64) (b0, b1, b2, a1, a2 float64) {
	w0 := 2 * math.Pi * cutoff / sampleRate
	A := math.Pow(10, gainDB/40)
	alpha := math.Sin(w0) / 2 * math.Sqrt((A+1/A)*(1/1.0-1)+2)

	switch filterType {
	case LowPass:
		alpha := math.Sin(w0) / 2
		b0 = (1 - math.Cos(w0)) / 2
		b1 = 1 - math.Cos(w0)
		b2 = (1 - math.Cos(w0)) / 2
		a0 := 1 + alpha
		a1 = -2 * math.Cos(w0)
		a2 = 1 - alpha
		b0 /= a0
		b1 /= a0
		b2 /= a0
		a1 /= a0
		a2 /= a0
	case HighPass:
		alpha := math.Sin(w0) / 2
		b0 = (1 + math.Cos(w0)) / 2
		b1 = -(1 + math.Cos(w0))
		b2 = (1 + math.Cos(w0)) / 2
		a0 := 1 + alpha
		a1 = -2 * math.Cos(w0)
		a2 = 1 - alpha
		b0 /= a0
		b1 /= a0
		b2 /= a0
		a1 /= a0
		a2 /= a0
	case LowShelf:
		b0 = A * ((A + 1) - (A-1)*math.Cos(w0) + 2*math.Sqrt(A)*alpha)
		b1 = 2 * A * ((A - 1) - (A+1)*math.Cos(w0))
		b2 = A * ((A + 1) - (A-1)*math.Cos(w0) - 2*math.Sqrt(A)*alpha)
		a0 := (A + 1) + (A-1)*math.Cos(w0) + 2*math.Sqrt(A)*alpha
		a1 = -2 * ((A - 1) + (A+1)*math.Cos(w0))
		a2 = (A + 1) + (A-1)*math.Cos(w0) - 2*math.Sqrt(A)*alpha
		b0 /= a0
		b1 /= a0
		b2 /= a0
		a1 /= a0
		a2 /= a0
	case HighShelf:
		b0 = A * ((A + 1) + (A-1)*math.Cos(w0) + 2*math.Sqrt(A)*alpha)
		b1 = -2 * A * ((A - 1) + (A+1)*math.Cos(w0))
		b2 = A * ((A + 1) + (A-1)*math.Cos(w0) - 2*math.Sqrt(A)*alpha)
		a0 := (A + 1) - (A-1)*math.Cos(w0) + 2*math.Sqrt(A)*alpha
		a1 = 2 * ((A - 1) - (A+1)*math.Cos(w0))
		a2 = (A + 1) - (A-1)*math.Cos(w0) - 2*math.Sqrt(A)*alpha
		b0 /= a0
		b1 /= a0
		b2 /= a0
		a1 /= a0
		a2 /= a0
	}

	return
}
