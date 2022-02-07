package main

import (
	"fmt"
	"time"
)

func main() {
	resultCh := make(chan int)
	fmt.Println("main started")
	go add(100, 200, resultCh)
	fmt.Println("Initiated the add operation")
	result := 0
	fmt.Println("result = ", result)
	fmt.Println("main completed")
}

func add(x, y int, resultCh chan int) {
	time.Sleep(4 * time.Second)
	result := x + y
	resultCh <- result
}
