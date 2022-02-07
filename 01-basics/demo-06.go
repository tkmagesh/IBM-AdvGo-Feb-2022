package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	resultCh := make(chan int)
	wg := &sync.WaitGroup{}
	fmt.Println("main started")
	wg.Add(1)
	go add(100, 200, resultCh, wg)
	fmt.Println("Initiated the add operation")
	wg.Wait() //result in a deadlock
	result := <-resultCh
	fmt.Println("result = ", result)
	fmt.Println("main completed")
}

func add(x, y int, resultCh chan int, wg *sync.WaitGroup) {
	time.Sleep(4 * time.Second)
	result := x + y
	resultCh <- result
	wg.Done()
}
