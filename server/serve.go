package main

import (
	"flag"
	"fmt"
	"os"

	// local
	"server/core"
	"server/file"
	"server/wiki"
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: serve -datapath=/path/to/data/folder")
	os.Exit(1)
}

func init() {
	var dataPrefix string
	defaultPath := ""
	flag.StringVar(&dataPrefix, "datapath", defaultPath, "the path to the data folder")
	flag.Parse()
	if dataPrefix == defaultPath {
		usage()
	}
	core.SetDataPrefix(dataPrefix)
}

func main() {
	wiki.New("wiki")
	wiki.New("wiki2") // because I can
	file.New("file")
	core.NewDefaultServer("wiki", false)
	core.StartServers(8080)
}
