package main

import (
	"fmt"
	"sync"
	"time"

	"context"
)

func main() {
	fmt.Println("Main started")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	go doSomething(ctx, wg)
	go func() {
		fmt.Println("hit ENTER to cancel")
		var input string
		fmt.Scanln(&input)
		cancel()
	}()
	wg.Wait()
	fmt.Println("Main completed")
}

func doSomething(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Canellation instruction received. returning")
			return
		default:
			time.Sleep(500 * time.Millisecond)
			fmt.Println("Doing something")
		}
	}
}
