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
	rootCtx, cancel := context.WithCancel(context.Background())
	valueCtx := context.WithValue(rootCtx, "root-key", "root-value")
	go doSomething(valueCtx, wg)
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
	//ctx3, cancel := context.WithCancel(ctx)
	ctx3 := context.WithValue(ctx, "key-3", "value-3")
	go doSomethingBoring(ctx3, wg)
	//fmt.Println("Data received through context ", ctx.Value("key-1"))
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
	fmt.Println("Data received through context for [root-key] = ", ctx.Value("root-key"))
	fmt.Println("Data received through context for [key-3] = ", ctx.Value("key-3"))
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
