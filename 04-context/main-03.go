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
	wg.Add(1)
	ctx1 := context.WithValue(ctx, "key-1", "value-1")
	go doSomethingElse(ctx1, wg)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("[doSomething] Canellation instruction received. returning")
			return
		default:
			time.Sleep(500 * time.Millisecond)
			fmt.Println("Doing something")
		}
	}
}

func doSomethingElse(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Data received through context ", ctx.Value("key-1"))
	for {
		select {
		case <-ctx.Done():
			fmt.Println("[doSomethingElse] Canellation instruction received. returning")
			return
		default:
			time.Sleep(250 * time.Millisecond)
			fmt.Println("Doing something else")
		}
	}
}
