package main

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/aaronellington/zendesk-go/examples/chat/gibsonTest/concur"
)

func main() {
	ctx := context.Background()

	c := runtime.NumGoroutine()
	fmt.Printf("There are %d goroutines", c)
	defer func() {
		c2 := runtime.NumGoroutine()
		if c != c2 {
			fmt.Printf("There are %d / %d goroutines", c, c2)
		}
	}()

	x := &fetchThing{}
	myThing := concur.New[[]byte](
		func() ([]byte, error) {
			return x.Fetch()
		},
		concur.Custom{
			Chan: time.NewTicker(time.Second * 5).C,
		},
	)

	// timeDelay := time.After(time.Second * 20)

	go myThing.Loop(ctx)
	for update := range myThing.Updates() {
		fmt.Println(string(update))
	}
}

type fetchThing struct{}

func (f *fetchThing) Fetch() ([]byte, error) {
	random := rand.Int()

	// Imagine that this is the net.Conn read io blocking action, that then returns a []byte
	return []byte(fmt.Sprintf("Hello, %d", random)), nil
}
