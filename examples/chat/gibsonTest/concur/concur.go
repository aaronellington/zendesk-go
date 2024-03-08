package concur

import (
	"context"
	"log"
	"time"
)

type FetchResult[T any] struct {
	Item T
	Err  error
}

func New[T any](fetcher func(context.Context) (T, error)) *Reader[T] {
	return &Reader[T]{
		fetcher: fetcher,
		closing: make(chan chan error),
		updates: make(chan FetchResult[T]),
	}
}

type Reader[T any] struct {
	fetcher func(context.Context) (T, error)
	closing chan bool
	updates chan FetchResult[T]
}

func (r *Reader[T]) Updates() <-chan FetchResult[T] {
	return r.updates
}

func (r *Reader[T]) Close() {
	// errc := make(chan error)
	r.closing <- true
	// return <-errc
}

//

func (r *Reader[T]) Loop(ctx context.Context) {
	var fetchDone chan FetchResult[T]
	//
	var queue []FetchResult[T]
	var err error

	counter := 0

	for {
		counter++
		log.Printf("Loop: %d\n", counter)

		var first FetchResult[T]
		var updates chan FetchResult[T]

		if len(queue) > 0 {
			first = queue[0]
			updates = r.updates
		}

		var startFetch <-chan time.Time

		if fetchDone == nil && err == nil {
			startFetch = time.After(0)
			// log.Println("signalling to start fetch...")
		}

		// Set up channels for cases
		select {
		case <-r.closing:
			log.Printf("!! -- This case was chosen for Loop %d, %s\n", counter, "closing")

			close(r.updates)
			return

		case <-startFetch:
			log.Printf("!! -- This case was chosen for Loop %d, %s\n", counter, "startFetch")
			fetchDone = make(chan FetchResult[T], 1)
			// log.Printf("starting fetch... - Fetchdone: %v\n", fetchDone)

			go func() {
				fetched, err := r.fetcher(ctx)
				fetchDone <- FetchResult[T]{fetched, err}
			}()

		case fetchResult := <-fetchDone:
			log.Printf("!! -- This case was chosen for Loop %d, %s\n", counter, "fetchdone")
			err = fetchResult.Err
			queue = append(queue, fetchResult)
			fetchDone = nil

		case updates <- first:
			log.Printf("!! -- This case was chosen for Loop %d, %s\n", counter, "updates")
			queue = queue[1:]
		}
	}
}
