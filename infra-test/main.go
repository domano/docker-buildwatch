package main

import "time"

// Just some long running code to test hot reloading

func main() {
	for {
		println("Running...")
		<-time.After(1 * time.Second)
	}
}
