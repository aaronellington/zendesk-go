package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"

	"time"

	"github.com/aaronellington/zendesk-go/examples/chat/gibsonTest/concur"
)

func sleep(seconds int, endSignal chan<- bool) {
	time.Sleep(time.Duration(seconds) * time.Second)
	endSignal <- true
}

func main() {
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
			time.Sleep(time.Nanosecond * time.Duration(rand.Int31()))

			return payload, nil
		},
	)

	time.AfterFunc(time.Second*1, func() {
		cancel()
	})

	go x.Loop(ctx)

	for update := range x.Updates() {
		if update.Err != nil {
			fmt.Println("err: ", update.Err)
			fmt.Println(x.Close())
			panic("er")
		}

		log.Printf("From the main: %s\n", update.Item)
	}

}
