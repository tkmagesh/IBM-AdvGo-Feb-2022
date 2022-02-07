package main

import (
	"fmt"
	"time"
)

func main() {
	c1 := time.After(3 * time.Second)
	c2 := time.After(5 * time.Second)
	c3 := time.After(1 * time.Second)

	for i := 0; i < 3; i++ {
		select {
		case <-c1:
			fmt.Println("c1 completed")
		case <-c2:
			fmt.Println("c2 completed")
		case <-c3:
			fmt.Println("c3 completed")
		default:
			fmt.Println("default triggered")
		}
	}
	fmt.Println("Done")
}
