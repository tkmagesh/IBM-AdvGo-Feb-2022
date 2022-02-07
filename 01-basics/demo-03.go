package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup
var tasks []func() = make([]func(), 20)

func main() {
	fmt.Println("main started")
	//create 20 tasks "dynamically"

	for i := 1; i <= 20; i++ {
		wg.Add(1)
		tasks[i-1] = func(idx int) func() {
			return func() {
				fmt.Println("Task ", idx, "invoked")
				wg.Done()
			}
		}(i)
	}
	fmt.Println("20 tasks are created")
	fmt.Println("tasks execution initiated")

	//initiate the execution of all the tasks
	for _, task := range tasks {
		go task()
	}

	//wait for the tasks to complete
	wg.Wait()
	fmt.Println("main completed")
}

//each task will print "Task <n> invoked"
