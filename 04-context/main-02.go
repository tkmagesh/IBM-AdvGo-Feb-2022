package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("Main started")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	done := make(chan bool)
	go doSomething(done, wg)
	go func() {
		fmt.Println("hit ENTER to cancel")
		var input string
		fmt.Scanln(&input)
		close(done)
	}()
	wg.Wait()
	fmt.Println("Main completed")
}

func doSomething(done chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	wg.Add(1)
	go doSomethingElse(done, wg)
	for {
		select {
		case <-done:
			fmt.Println("[doSomething] Canellation instruction received. returning")
			return
		default:
			time.Sleep(500 * time.Millisecond)
			fmt.Println("Doing something")
		}
	}
}

func doSomethingElse(done chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-done:
			fmt.Println("[doSomethingElse] Canellation instruction received. returning")
			return
		default:
			time.Sleep(500 * time.Millisecond)
			fmt.Println("Doing something else")
		}
	}
}
