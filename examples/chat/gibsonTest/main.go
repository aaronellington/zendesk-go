package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"runtime"

	"time"
)

func sleep(seconds int, endSignal chan<- bool) {
	time.Sleep(time.Duration(seconds) * time.Second)
	endSignal <- true
}

func main() {
	c := runtime.NumGoroutine()
	fmt.Printf("There are %d goroutines", c)
	defer func() {
		c2 := runtime.NumGoroutine()
		// if c != c2 {
		fmt.Printf("There are %d / %d goroutines", c, c2)
		// }
	}()
	myThing := reader{
		fetcher: &fetchThing{},
		closing: make(chan chan error),
		updates: make(chan []byte),
	}

	// for {
	// 	fmt.Println(rand.Int())
	// }

	// timeDelay := time.After(time.Second * 20)

	go myThing.loop(context.Background())
	for update := range myThing.updates {
		select {
		// case <-timeDelay:
		// 	panic("aa")
		default:
			fmt.Println(string(update))
		}
	}
}

type Fetcher interface {
	Fetch() (frames []byte, err error)
}

type fetchThing struct{}

func (f *fetchThing) Fetch() ([]byte, error) {
	random := rand.Int()

	// Imagine that this is the net.Conn read io blocking action, that then returns a []byte
	return []byte(fmt.Sprintf("Hello, %d", random)), nil
}

// type ReadLooper interface {
// 	Updates() <-chan []byte
// 	Close() error
// }

type reader struct {
	fetcher Fetcher
	closing chan chan error // This is a request / response structure
	updates chan []byte
}

func (r *reader) Updates() <-chan []byte {
	return r.updates
}

func (r *reader) Close() error {
	errc := make(chan error)
	r.closing <- errc
	return <-errc
}

type fetchResult struct {
	item []byte
	err  error
}

func (r *reader) loop(ctx context.Context) {
	// Mutable state
	// endSignal := make(chan bool)
	// go sleep(20, endSignal)

	var fetchDone chan fetchResult
	//
	var pendingMessage [][]byte
	var err error
	counter := 0

	pinger := time.NewTicker(time.Second * 5)
	defer pinger.Stop()

	for {
		counter++
		log.Printf("Loop: %d\n", counter)
		var first []byte
		var updates chan []byte

		if len(pendingMessage) > 0 {
			first = pendingMessage[0]
			updates = r.updates
		}

		var startFetch <-chan time.Time

		if fetchDone == nil {
			startFetch = time.After(time.Second * 1)
			log.Println("signalling to start fetch...")
		}

		// Set up channels for cases
		select {
		case <-pinger.C:
			log.Println("I've pinged!")
			// case <-endSignal:
		// 	close(r.updates)
		// 	// panic("Show me the stacks!")
		case <-ctx.Done():
			// errc <- context.Cause("")
			close(r.updates)
			return

		case errc := <-r.closing:
			errc <- err
			close(r.updates)
			return

		case <-startFetch:
			fetchDone = make(chan fetchResult, 1)
			log.Println("starting fetch...")

			go func() {
				fetched, err := r.fetcher.Fetch()
				time.Sleep(time.Second)
				fetchDone <- fetchResult{fetched, err}
			}()

		case fetchedItems := <-fetchDone:
			log.Println("fetch complete...")

			if fetchedItems.err != nil {
				break
			}

			pendingMessage = append(pendingMessage, fetchedItems.item)
			fetchDone = nil

		case updates <- first:
			pendingMessage = pendingMessage[1:]
		}
	}
}
