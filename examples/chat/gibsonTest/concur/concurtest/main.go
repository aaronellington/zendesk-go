package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"

	"time"

	"github.com/aaronellington/zendesk-go/examples/chat/gibsonTest/concur"
)

func sleep(seconds int, endSignal chan<- bool) {
	time.Sleep(time.Duration(seconds) * time.Second)
	endSignal <- true
}

func main() {

	defer func() {
		time.Sleep(time.Millisecond * 500)
		c := runtime.NumGoroutine()

		if c > 1 {

			fmt.Printf("There are %d extra goroutines compared to beginning of run", c-1)
		}
	}()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	payloads := []string{
		"Hello",
		"There",
		"Said",
		"General",
		"Kenobi",
	}

	x := concur.New[string](
		func(ctx context.Context) (string, error) {
			if len(payloads) == 0 {
				return "", errors.New("no more data")
			}
			payload := payloads[0]
			payloads = payloads[1:]
			// time.Sleep(time.Nanosecond * time.Duration(rand.Int31()))

			return payload, nil
		},
	)

	time.AfterFunc(time.Second*1, func() {
		cancel()
	})

	go x.Loop(ctx)

	// for update := range x.Updates() {
	// 	if update.Err != nil {
	// 		fmt.Println("err: ", update.Err)
	// 		x.Close()
	// 		panic("er")
	// 	}

	// 	log.Printf("From the main: %s\n", update.Item)
	// }
	// go func() {
	// 	time.Sleep(time.Second * 100)
	// }()

	for {
		select {
		case update := <-x.Updates():
			if update.Err != nil {
				fmt.Println("err: ", update.Err)
				x.Close()
				panic("er")
			}

			log.Printf("From the main: %s\n", update.Item)

		case <-ctx.Done():
			log.Printf("context reason: %s\n", context.Cause(ctx))
			x.Close()
			panic("er")

		}

	}

}
