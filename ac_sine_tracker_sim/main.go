package main

import "fmt"

func main() {
	ch := MakeSource(60)
	for i := 0; ; i++ {
		x := <-ch
		if i%1000 == 0 {
			fmt.Println(x)
		}
	}
}
