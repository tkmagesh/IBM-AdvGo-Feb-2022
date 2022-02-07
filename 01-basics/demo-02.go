package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {
	fmt.Println("main started")
	wg.Add(1)
	go f1()
	f2()
	wg.Wait()
	fmt.Println("main completed")
}

func f1() {
	fmt.Println("f1 is invoked")
	wg.Done()
}

func f2() {
	fmt.Println("f2 is invoked")
}
