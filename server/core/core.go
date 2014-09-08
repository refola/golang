package core

import (
	"fmt"
	"runtime"
	"time"
)

func init() {
	go func() {
		var m *runtime.MemStats
		p := func(name string, what interface{}) {
			fmt.Printf("\t%s: %v\n", name, what)
		}
		printStats := func() {
			runtime.ReadMemStats(m)
			p("Alloc", m.Alloc)
			p("TotalAlloc", m.TotalAlloc)
		}
		for {
			time.Sleep(10 * 1e9)
			fmt.Println("Calling garbage collector.")
			fmt.Println("Stats before:")
			printStats()
			runtime.GC()
			fmt.Println("Stats after:")
			printStats()
		}
	}()
}
