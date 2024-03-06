package concur

import (
	"context"
	"log"
	"time"
)

type fetchResult[Y any] struct {
	item Y
	err  error
}

func New[T any](fetcher func() (T, error)) *Reader[T] {
	return &Reader[T]{
		fetcher: fetcher,
		closing: make(chan chan error),
		updates: make(chan T),
	}
}

type Reader[T any] struct {
	fetcher func() (T, error)
	closing chan chan error
	updates chan T
}

func (r *Reader[T]) Updates() <-chan T {
	return r.updates
}

func (r *Reader[T]) Close() error {
	errc := make(chan error)
	r.closing <- errc
	return <-errc
}

func (r *Reader[T]) Loop(ctx context.Context) {
	// Mutable state
	// endSignal := make(chan bool)
	// go sleep(20, endSignal)

	var fetchDone chan fetchResult[T]
	//
	var queue []T
	var err error
	counter := 0

	for {
		counter++
		log.Printf("Loop: %d\n", counter)
		var first T
		var updates chan T

		if len(queue) > 0 {
			first = queue[0]
			updates = r.updates
		}

		var startFetch <-chan time.Time

		if fetchDone == nil {
			startFetch = time.After(time.Second * 1)
			log.Println("signalling to start fetch...")
		}

		// Set up channels for cases
		select {
		case <-ctx.Done():
			close(r.updates)
			return

		case errc := <-r.closing:
			errc <- err
			close(r.updates)
			return

		case <-startFetch:
			fetchDone = make(chan fetchResult[T], 1)
			log.Println("starting fetch...")

			go func() {
				fetched, err := r.fetcher()
				time.Sleep(time.Second)
				fetchDone <- fetchResult[T]{fetched, err}
			}()

		case fetchedItems := <-fetchDone:
			log.Println("fetch complete...")

			if fetchedItems.err != nil {
				break
			}

			queue = append(queue, fetchedItems.item)
			fetchDone = nil

		case updates <- first:
			queue = queue[1:]
		}
	}
}
