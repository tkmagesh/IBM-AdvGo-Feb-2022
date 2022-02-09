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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	go doSomething(ctx, wg)
	wg.Wait()
	fmt.Println("Main completed")
}

func doSomething(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	wg.Add(1)
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	go doSomethingElse(ctx2, wg)
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
	wg.Add(1)
	ctx3, cancel := context.WithCancel(ctx)
	defer cancel()
	go doSomethingBoring(ctx3, wg)
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

func doSomethingBoring(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Data received through context ", ctx.Value("key-1"))
	for {
		select {
		case <-ctx.Done():
			fmt.Println("[doSomethingBoring] Canellation instruction received. returning")
			return
		default:
			time.Sleep(100 * time.Millisecond)
			fmt.Println("Doing something boring")
		}
	}
}
