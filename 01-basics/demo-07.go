package main

import (
	"fmt"
	"time"
)

func main() {
	resultCh := make(chan int, 1)
	fmt.Println("main started")
	go add(100, 200, resultCh)
	fmt.Println("Initiated the add operation")
	time.Sleep(8 * time.Second)
	fmt.Println("8 seconds elapsed")
	result := <-resultCh
	fmt.Println("result = ", result)
	fmt.Println("main completed")
}

func add(x, y int, resultCh chan int) {
	time.Sleep(4 * time.Second)
	result := x + y
	resultCh <- result
	fmt.Println("Write succeeded")
}
