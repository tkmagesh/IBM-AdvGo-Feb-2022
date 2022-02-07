/* solving demo-03.go using channels */
package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup
var tasks = make(chan func(), 20)

func main() {
	fmt.Println("main started")
	//create 20 tasks "dynamically"

	for i := 1; i <= 20; i++ {
		wg.Add(1)
		tasks <- func(idx int) func() {
			return func() {
				fmt.Println("Task ", idx, "invoked")
				wg.Done()
			}
		}(i)
	}
	close(tasks)
	fmt.Println("20 tasks are created")
	fmt.Println("tasks execution initiated")

	//initiate the execution of all the tasks
	for task := range tasks {
		go task()
	}

	//wait for the tasks to complete
	wg.Wait()
	fmt.Println("main completed")
}

//each task will print "Task <n> invoked"
