// Get input from the human-reading telemach.
// TODO: Add stuff for making this work (seamlessly) with other inputs like via emokit

package telemach

import (
	pa "code.google.com/p/portaudio-go/portaudio"
	//"errors"
	"fmt"
	//"os"
)

func init() {
	err := pa.Initialize() // seems to be needed before doing other stuff with pa
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		// os.Exit(1)
	}
	// defer pa.Terminate()
}

// Input should resolve to arrays of something to represent amplitudes.
type TMSnippet struct {
	SampleHz float64   // How often the input is sampled
	Samples  []float32 // Slice containing a bunch of samples in whatever form is correct (TODO: update type from int to whatever is needed)
}

// Wrapper for a PortAudio device, keeping the details out of the way.
type Device struct {
	paDeviceInfo *pa.DeviceInfo
	Name         string
}

// Get PortAudio device list and wrap them into this package's Device type.
func ListDevices() ([]Device, error) {
	paDevices, err := pa.Devices()
	if err != nil {
		return nil, err
	}
	devices := make([]Device, len(paDevices))
	for i, v := range paDevices {
		devices[i].paDeviceInfo = v
		devices[i].Name = fmt.Sprintf("%v %v", i, v.Name)
	}
	return devices, nil
}

// Get a channel to read convenient data snippets from a device.
//  returns
//	<-chan TMSnippet	for caller to get audio samples
//	chan<- bool		for caller to send on when wanting to gracefully shutdown the audio stream and associated chans (true and false are treated the same here)
//	<-chan error		for the data-sending goroutine to report any errors
//	error			nil unless the stream can't be opened
// Note that this function doesn't (yet) automatically close if errors happen. It is currently the caller's responsibility to check if the error can be ignored and if not, to send something to the chan<- bool.
func MakeInputChannel(device Device) (<-chan *TMSnippet, chan<- bool, <-chan error, error) {
	streamParams := pa.LowLatencyParameters(device.paDeviceInfo, nil)
	streamParams.FramesPerBuffer = 64 // Number taken from portaudio-go/portaudio/examples/record.go line 70
	fmt.Printf("Stream parameters: %v\n", streamParams)
	buffer := make([]float32, streamParams.FramesPerBuffer)
	stream, err := pa.OpenStream(streamParams, buffer)
	if err != nil {
		return nil, nil, nil, err
	}
	tmschan := make(chan *TMSnippet, 10)
	closechan := make(chan bool)
	errchan := make(chan error)
	go func() {
		if err := stream.Start(); err != nil {
			errchan <- err
		}
		defer func() {
			err := stream.Close()
			if err != nil {
				errchan <- err
			}
			close(closechan)
			close(errchan)
		}()
		for {
			err := stream.Read()
			if err != nil {
				errchan <- err
				break
			}
			tms := new(TMSnippet)
			tms.SampleHz = streamParams.SampleRate
			tms.Samples = make([]float32, len(buffer))
			copy(tms.Samples, buffer)
			tmschan <- tms
		}
		/*
			for {
				if _, end := <-closechan; end {
					break
				}
				var tms TMSnippet
				tms.SampleHz = streamParams.SampleRate
				if framesAvailable, err := stream.AvailableToRead(); err == nil {
					if framesAvailable <= 0 {
						errchan <- errors.New("No frames available.")
						continue
					}
					// TODO: per TODO below, see if len(buffer) is more appropriate than framesAvailable
					tms.Samples = make([]float32, framesAvailable)
					// TODO: see if this gets stream.AvailableToRead()'s "maximum number of frames that can be read from the stream without blocking or busy waiting" described in portaudio.h or if it gets len(buffer) frames per method comment
					err = stream.Read()
					if err != nil {
						errchan <- err
					} else {
						n := copy(tms.Samples, buffer)
						if n != framesAvailable {
							errchan <- errors.New("Didn't get expected number of frames.")
						}
						tmschan <- &tms
					}
				} else {
					errchan <- err
				}
			}
		*/
	}()
	return tmschan, closechan, errchan, nil
}
