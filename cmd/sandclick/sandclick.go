// sandclick.go - automate the xkcd Time sandcastle builder at http://castle.chirpingmustard.com

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
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
}

// get the mNP progress of the current newpix,
func currentmNP(currentpix int64) int64 {
	newpix := newpix(currentpix)

	// current mNP = 1000 * nanos into current NP / total nanos in current NP
	npNanos := newpix / time.Nanosecond
	return 1000 * (time.Now().UnixNano() % int64(npNanos)) / int64(npNanos)
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

// Runs cmd.Run() repeatedly after true is received on status and stops when false is received. Stops running when channel is closed.
func cmdRepeater(status <-chan bool, cmdBuilder func() *exec.Cmd) {
	defer func() { recover(); fmt.Println("channel closed") }()
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
	sleepTomNP(curpix, wakemNP)
	ch <- true
	sleepTomNP(curpix, sleepmNP)
	ch <- false
}

func main() {
	// get all the vars
	numNewpix := flag.Int64("numnewpix", -1, "Mandatory -- The number of NewPix to click for, 0 for infinite, negative for displaying help.")
	xserver := flag.String("xserver", os.Getenv("DISPLAY"), "Which X server to send the clicks to (defaults to the one this command is ran from).")
	startpix := flag.Int64("currentpix", 1, "Which pix the game is on.")
	npbotmNP := flag.Int64("npbotmnp", 400, "How many mNP the NewPixBots need to avoid being ninja'd (0 if there's no ninja).")
	helpwanted := flag.Bool("help", false, "Print the help message.")
	flag.Parse()

	// verify parameter validity and display help if anything's invalid (or if help's requested)
	if *helpwanted || *numNewpix < 0 || *startpix < 1 || *npbotmNP < 0 {
		usage()
		os.Exit(1)
	}

	// make clicker
	ch := make(chan bool, 1)
	go cmdRepeater(ch, func() *exec.Cmd { return clickNCmd(5, *xserver) })

	// finish current newpix
	curpix := *startpix
	var margin int64 = 3 // how many extra mNP on either side of newpixbot ninja'ing time to wait
	wakemNP := *npbotmNP + margin
	sleepmNP := 1000 - margin
	if mNP := currentmNP(curpix); mNP < sleepmNP {
		clickLoopIteration(curpix, wakemNP, sleepmNP, ch)
	}
	for curpix++; curpix < *startpix+*numNewpix; curpix++ {
		clickLoopIteration(curpix, wakemNP, sleepmNP, ch)
	}
}
