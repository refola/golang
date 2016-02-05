// sandclick.go - automate the xkcd Time sandcastle builder at http://castle.chirpingmustard.com

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// convert times to Outside time
const longpix = time.Hour
const shortpix = longpix / 2 // shortpix are half as long in Outside time

func newpix(currentpix int64) time.Duration {
	if currentpix <= 240 { // newpix up to 240 are shortpix, per http://xkcd-time.wikia.com/wiki/WarpPix
		return shortpix
	} else {
		return longpix
	}
	// Debian's golang package's compiler isn't updated enough to realize that this will never be reached.
	panic("Unreachable!")
}

// get the mNP progress of the current newpix,
func currentmNP(currentpix int64) int64 {
	newpix := newpix(currentpix)

	// current mNP = 1000 * nanos into current NP / total nanos in current NP
	npNanos := int64(newpix / time.Nanosecond)
	return (1000 * (time.Now().UnixNano() % npNanos)) / npNanos
}

func sleepTomNP(curpix, mNPToSleepTo int64) {
	mNP := newpix(curpix) / 1000
	sleepmNP := mNPToSleepTo - currentmNP(curpix)
	if sleepmNP > 0 {
		fmt.Printf("Sleeping for %vmNP.\n", sleepmNP)
		time.Sleep(mNP * time.Duration(sleepmNP))
	}
}

// build the command to click n times
func clickNCmd(n int, server string) *exec.Cmd {
	args := make([]string, 2+n)
	args[0] = "-x"
	args[1] = server
	for i := 2; i <= n; i++ {
		args[i] = "mouseclick 1"
	}
	return exec.Command("xte", args...)
}

// Runs cmd.Run() repeatedly with a given delay between runnings after true is received on status and stops when false is received. Stops running when channel is closed.
func cmdRepeater(delay time.Duration, status <-chan bool, cmdBuilder func() *exec.Cmd) {
	defer func() { recover(); fmt.Println("cmdRepeater exiting: channel closed") }()
	active := <-status
	for {
		if len(status) > 0 {
			active = <-status
		}
		if active {
			err := cmdBuilder().Run()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			time.Sleep(delay)
		} else {
			time.Sleep(time.Second)
		}
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "sandclick: xkcd Time Sandcastle Builder autoclicker for http://castle.chirpingmustard.com/")
	flag.PrintDefaults()
}

func clickLoopIteration(curpix, wakemNP, sleepmNP int64, ch chan<- bool) {
	fmt.Printf("Now on newpix number %v.\n", curpix)
	sleepTomNP(curpix, wakemNP)
	fmt.Println("Activating clicking.")
	ch <- true
	sleepTomNP(curpix, sleepmNP)
	fmt.Println("Deactivating clicking.")
	ch <- false
	fmt.Println("Waiting for next newpix.")
	sleepTomNP(curpix, 1001)
}

func main() {
	// So the command-running (clicking) goroutine can be ran separately....
	procs := 2
	fmt.Printf("Switching from using %v cores to %v.\n", runtime.GOMAXPROCS(procs), procs)

	// get all the vars
	numNewpix := flag.Int64("numnewpix", -1, "Mandatory -- The number of NewPix to click for, 0 for infinite, negative for displaying help.")
	xserver := flag.String("xserver", os.Getenv("DISPLAY"), "Which X server to send the clicks to (defaults to the one this command is ran from).")
	startpix := flag.Int64("currentpix", 1, "Which pix the game is on.")
	npbotmNP := flag.Int64("npbotmnp", 400, "How many mNP the NewPixBots need to avoid being ninja'd (0 if there's no ninja).")
	helpwanted := flag.Bool("help", false, "Print the help message.")
	delayms := flag.Int("delayms", 5, "How many milliseconds to delay successive click commands.")
	flag.Parse()

	// verify parameter validity and display help if anything's invalid (or if help's requested)
	if *helpwanted || *numNewpix < 0 || *startpix < 1 || *npbotmNP < 0 {
		usage()
		os.Exit(1)
	}

	// make clicker
	ch := make(chan bool, 1)
	go cmdRepeater(time.Millisecond*time.Duration(*delayms), ch, func() *exec.Cmd { return clickNCmd(5, *xserver) })

	// finish current newpix
	curpix := *startpix
	var margin int64 = 3 // how many extra mNP on either side of newpixbot ninja'ing time to wait
	wakemNP := *npbotmNP + margin
	sleepmNP := 1000 - margin
	if mNP := currentmNP(curpix); mNP < sleepmNP {
		clickLoopIteration(curpix, wakemNP, sleepmNP, ch)
	}

	// do other newpix
	stopAt := *startpix + *numNewpix
	if *numNewpix == 0 {
		stopAt = 1000000 // At one hour per newpix, this should take over 100 years to stop.
	}
	for curpix++; curpix < stopAt; curpix++ {
		clickLoopIteration(curpix, wakemNP, sleepmNP, ch)
	}
}
