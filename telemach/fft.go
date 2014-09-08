// Do a Fast Fourier Transform on some telepathy machine data for further processing.

package telemach

import (
	"github.com/mjibson/go-dsp/fft"
)

func init() {
	fft.SetWorkerPoolSize(1) // Each stream needs its own core. TODO: See if this actually works.... I'm trying to make it use one distinct single-threaded worker per channel.
}

type Frequencies struct {
	DeltaHz float64   // how much the frequency changes each index
	Freqs   []float64 // slice containing the actual frequencies
}

func cplxsToReals(cplxs []complex128) []float64 {
	reals := make([]float64, len(cplxs))
	for i, v := range cplxs {
		reals[i] = real(v)
	}
	return reals
}

func SamplesToFrequencies(samples *TMSnippet) Frequencies {
	var freqs Frequencies
	f64samples := make([]float64, len(samples.Samples))
	for i, v := range samples.Samples {
		f64samples[i] = float64(v)
	}
	freqs.Freqs = cplxsToReals(fft.FFTReal(f64samples))
	freqs.DeltaHz = 1 / samples.SampleHz // TODO: check if this even makes sense
	return freqs
}

func FrequenciesToSamples(freqs Frequencies) *TMSnippet {
	var samples TMSnippet
	f64samples := cplxsToReals(fft.IFFTReal(freqs.Freqs))
	samples.Samples = make([]float32, len(f64samples))
	for i, v := range f64samples {
		samples.Samples[i] = float32(v)
	}
	samples.SampleHz = 1 / freqs.DeltaHz // TODO: check if this even makes sense
	return &samples
}
