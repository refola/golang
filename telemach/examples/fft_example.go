package main

import (
	tm "code.google.com/p/refola/telemach"
	"fmt"
	"math"
)

func sinePeriod(samplesPerPeriod int) []float64 {
	delta := 2 * math.Pi / float64(samplesPerPeriod)
	samples := make([]float64, samplesPerPeriod)
	for i, _ := range samples {
		samples[i] = math.Sin(float64(i) * delta)
	}
	return samples
}

func repeat(arr []float64, times int) []float64 {
	ret := make([]float64, len(arr)*times)
	for i, _ := range ret {
		ret[i] = arr[i%len(arr)]
	}
	return ret
}

// make a uniform sine wave with given number of crests and total samples
func sinArr(crests, totalSamples int) []float64 {
	return repeat(sinePeriod(totalSamples/crests), crests+1)[:totalSamples]
}

// return the indexes of arr that correspond to the top n values in arr
func topN(arr []float64, n int) []int {
	ret := make([]int, n)
	tops := make(map[int]float64)
	// TODO
}

func main() {
	var tms tm.TMSnippet
	tms.SampleHz = 1 // Only kept track of by FFT functions, not used for actual FFT
	s := sinArr(5, 100)
	tms.Samples = make([]float32, len(s))
	for i, v := range s {
		tms.Samples[i] = float32(v)
	}
	fmt.Printf("Samples: %v\n", (tms.Samples))
	freqs := tm.SamplesToFrequencies(&tms)
	fmt.Printf("Frequencies: %v\n", (freqs.Freqs))
}
