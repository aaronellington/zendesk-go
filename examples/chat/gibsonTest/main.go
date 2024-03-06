package main

import (
	"fmt"
	"log"
	"time"
)

func sleep(seconds int, endSignal chan<- bool) {
	time.Sleep(time.Duration(seconds) * time.Second)
	endSignal <- true
}

func main() {
	endSignal := make(chan bool, 1)
	go sleep(8, endSignal)

	var next time.Time
	next = time.Now().Add(time.Second * 15)
	for {
		var delay time.Duration
		if now := time.Now(); next.After(now) {
			delay = next.Sub(now)
		}
		startFetch := time.After(delay)
		log.Println("I am not sleepy, because I got this message!")

		select {
		case <-endSignal:
			fmt.Println("The end!")
		case <-startFetch:
			next = time.Now().Add(time.Second * 2)
			fmt.Println("There's no more time to this. Exiting!")

		}
	}
}
