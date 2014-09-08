// Generate noisy waves for the tracker
package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Adds randomness to the amplitudes passed through
func ampNoise(rx <-chan float64, tx chan<- float64) {
	for {
		tx <- (rand.Float64()/10 - 0.05 + <-rx)
	}
}

// Inserts or removes data points to simulate frequency changes
func skewNoise(rx <-chan float64, tx chan<- float64) {
	for {
		r := rand.Float32()
		in := <-rx
		switch {
		case r < 0.1:
		case r >= 0.03 && r < 0.2:
			tx <- in
			tx <- in
		default:
			tx <- in
		}
	}
}

// Sends a sine wave of amplitude 1 with given frequency and granularity
// ch: The channel this will send values on
// freq: The frequency in hertz to oscillate the wave with
// grains: The number of samples each wave should have
func oscillate(ch chan<- float64, freq float64, grains int) {
	cycle := 1 / freq
	tick := cycle / float64(grains) // seconds between number generations
	nanos := int64(tick * float64(time.Second/time.Nanosecond))
	dur := time.Nanosecond * time.Duration(nanos)
	fmt.Println("The time between oscillations is: " + dur.String())
	for {
		for count := float64(0); count <= cycle; count += tick {
			ch <- math.Sin(count * 2 * math.Pi / cycle)
			time.Sleep(dur)
		}
	}
}

// Gives a channel to listen to a noisy sine wave with
func MakeSource(Hz float64) <-chan float64 {
	oscch := make(chan float64)
	go oscillate(oscch, Hz, 10000)
	fromamp := make(chan float64)
	go ampNoise(oscch, fromamp)
	fromskew := make(chan float64)
	go skewNoise(fromamp, fromskew)
	return fromskew
}
