/* solving demo-03.go using channels */
package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup
var tasks = make(chan func())

func main() {
	fmt.Println("main started")
	//create 20 tasks "dynamically"
	wg.Add(20)
	go func() {
		for i := 1; i <= 20; i++ {
			tasks <- func(idx int) func() {
				return func() {
					fmt.Println("Task ", idx, "invoked")
					wg.Done()
				}
			}(i)
			fmt.Println("Task ", i, "created")
		}
		close(tasks)
		fmt.Println("20 tasks are created")
	}()

	go func() {
		fmt.Println("tasks execution initiated")

		//initiate the execution of all the tasks
		for task := range tasks {
			go task()
		}
	}()

	//wait for the tasks to complete
	wg.Wait()
	fmt.Println("main completed")
}

//each task will print "Task <n> invoked"
