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
		done <- true
	}()
	wg.Wait()
	fmt.Println("Main completed")
}

func doSomething(done chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-done:
			fmt.Println("Canellation instruction received. returning")
			return
		default:
			time.Sleep(500 * time.Millisecond)
			fmt.Println("Doing something")
		}
	}
}
