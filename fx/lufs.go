package fx

import (
	"math"
)

const (
	// shelving filter coefficiaent
	a1 = -1.69065929318241 //(-)
	a2 = 0.73248077421585
	b0 = 1.53512485958697
	b1 = -2.69169618940638 // (-)
	b2 = 1.19839281085285
)

const (
	hpCoef_a1 = -1.99004745483398 // (-)
	hpCoef_a2 = 0.99007225036621
	hpCoef_b0 = 1.0
	hpCoef_b1 = -2.0 //(-)
	hpCoef_b2 = 1.0
)

// func CalculateLufs(sampleRate int, sampleF []float64) float64 {
// 	// //	// DSP filter
// 	dt := 1. / float64(sampleRate) //sample rate
// 	wc := 2 * math.Pi * 60 * dt    //60 hz
// 	hp := butter.NewHighPass2(wc)
// 	hs := NewIIRFilter(b0, b1, b2, a1, a2)
// 	// high Shelving
// 	for i, v := range sampleF {
// 		sampleF[i] = hs.Process(v)
// 	}
// 	//high pass
// 	hp.NextS(sampleF, sampleF)
// 	// sampleF = applyKWeighting(sampleF)
// 	sampleF = calculateRMS(sampleF, sampleRate)
// 	sampleF = applyGating(sampleF)
// 	return calculateLUFS(sampleF)
// }

// func applyPreFiltering(samples []float64) []float64 {
// 	// Apply a high-pass filter with a cutoff frequency of 30 Hz
// 	// Simplified Butterworth filter design for demonstration
// 	const a1 = -1.995 * 0.999
// 	const a2 = 0.995 * 0.999
// 	const b0 = 1.0
// 	const b1 = -2.0
// 	const b2 = 1.0

// 	x1, x2 := 0.0, 0.0
// 	y1, y2 := 0.0, 0.0

// 	filtered := make([]float64, len(samples))
// 	for i, x0 := range samples {
// 		y0 := b0*x0 + b1*x1 + b2*x2 - a1*y1 - a2*y2
// 		filtered[i] = y0

// 		x2 = x1
// 		x1 = x0
// 		y2 = y1
// 		y1 = y0
// 	}

// 	return filtered
// }

// func applyKWeighting(samples []float64) []float64 {
// 	// Apply a shelf filter to emulate K-weighting
// 	const a1 = -1.69065929318241
// 	const a2 = 0.73248077421585
// 	const b0 = 1.53512485958697
// 	const b1 = -2.69169618940638
// 	const b2 = 1.19839281085285

// 	x1, x2 := 0.0, 0.0
// 	y1, y2 := 0.0, 0.0

// 	filtered := make([]float64, len(samples))
// 	for i, x0 := range samples {
// 		y0 := b0*x0 + b1*x1 + b2*x2 - a1*y1 - a2*y2
// 		filtered[i] = y0

// 		x2 = x1
// 		x1 = x0
// 		y2 = y1
// 		y1 = y0
// 	}

// 	return filtered
// }

///////////////////

func LUFS(samples []float64, sampleRate int) float64 {
	shelv := NewIIRFilter(ComputeCoefficients(HighShelf, 2000.0, 48000.0, 4))
	hp := NewIIRFilter(ComputeCoefficients(HighPass, 30, 48000.0, 0))
	out := make([]float64, len(samples))

	for i, s := range samples {
		out[i] = shelv.Process(s)
	}

	for i, s := range out {
		out[i] = hp.Process(s)
	}

	out = calculateRMS(out, sampleRate)
	out = applyGating(out)
	return calculateLUFS(out)

}

func calculateLUFS(rmsValues []float64) float64 {
	sum := 0.0
	for _, rms := range rmsValues {
		sum += rms * rms
	}
	meanSquare := sum / float64(len(rmsValues))
	lufs := -0.691 + (10 * math.Log10(meanSquare))
	return lufs
}

func applyGating(rmsValues []float64) []float64 {
	// Implement gating according to ITU-R BS.1770
	// Simplified here for demonstration
	threshold := -70.0
	gatedRMSValues := make([]float64, 0, len(rmsValues))
	for _, rms := range rmsValues {
		if 20*math.Log10(rms) > threshold {
			gatedRMSValues = append(gatedRMSValues, rms)
		}
	}
	return gatedRMSValues
}

func calculateRMS(samples []float64, windowSize int) []float64 {
	rmsValues := make([]float64, 0, len(samples)/windowSize)
	for i := 0; i < len(samples); i += windowSize {
		sum := 0.0
		for j := 0; j < windowSize && i+j < len(samples); j++ {
			sum += samples[i+j] * samples[i+j]
		}
		rms := math.Sqrt(sum / float64(windowSize))
		rmsValues = append(rmsValues, rms)
	}
	return rmsValues
}
