/*
Tweak input from telepathy machine and send something useful to output.

Needs:
 * Some way to get data from telepathy machine (currently obtained via cheap USB sound card)
 * Processing to get rid of as much noise as possible (especially 60 Hz ambient from local power infrastructure) and amplify the remaining signal (FFT from input, remove noise by frequency, FFT back to output?)
 * Some way to send the noise-reduced signal to the telepathy machine's output

*/

package main

import (
	tm "code.google.com/p/refola/telemach"
	"errors"
	"fmt"
	"os"
)

func showDevices(devs []tm.Device) {
	fmt.Println("Here are the devices.")
	for _, v := range devs {
		fmt.Println(v.Name)
	}
}

func selectDevice(devs []tm.Device) (tm.Device, error) {
	fmt.Println("Please enter a device number.")
	var input int
	n, err := fmt.Scanln(&input)
	if n != 1 || err != nil {
		if err == nil {
			err = errors.New("User entered more than one thing.")
		}
		return tm.Device{}, err
	}
	if input < 0 || input > len(devs) {
		return tm.Device{}, errors.New("User entered an invalid number.")
	}
	return devs[input], nil
}

func main() {
	var err error
	e := func() {
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
	devs, err := tm.ListDevices()
	e()
	showDevices(devs)
	dev, err := selectDevice(devs)
	e()
	fmt.Printf("Device selected: %v\n", dev.Name)
	datachan, closechan, errchan, err := tm.MakeInputChannel(dev)
	e()
	fmt.Println("Got data channel, etc....")
	go func() {
		for err := range errchan {
			fmt.Printf("Error from errchan: $v\n", err)
		}
	}()
	fmt.Println("Attempting to read from data channel....")
	tms := <-datachan
	fmt.Println("Got first sample!")
	go func() { closechan <- true }()
	fmt.Printf("We have data at %v Hz! here it is: %v\n", tms.SampleHz, tms.Samples)
	ffted := tm.SamplesToFrequencies(tms)
	fmt.Printf("We have a Fast Fourier Transformed sample with consecutive frequencies varying by %v Hz! Here are the frequencies: %v\n", ffted.DeltaHz, ffted.Freqs)
	unffted := tm.FrequenciesToSamples(ffted)
	fmt.Printf("Back from FFT! We have data at %v Hz! here it is: %v\n", unffted.SampleHz, unffted.Samples)
}
