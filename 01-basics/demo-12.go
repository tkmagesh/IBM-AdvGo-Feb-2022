/* mutex */
package main

import (
	"fmt"
	"sync"
)

type MyWaitGroup struct {
	counter int
	sync.Mutex
	done chan bool
}

func NewMyWaitGroup() *MyWaitGroup {
	return &MyWaitGroup{
		done: make(chan bool, 1),
	}
}

func (mwg *MyWaitGroup) Add(count int) {
	mwg.Lock()
	{
		mwg.counter += count
	}
	mwg.Unlock()
}

func (mwg *MyWaitGroup) Done() {
	mwg.Lock()
	{
		mwg.counter -= 1
		if mwg.counter == 0 {
			mwg.done <- true
		}
	}
	mwg.Unlock()
}

func (mwg *MyWaitGroup) Wait() {
	for {
		select {
		case <-mwg.done:
			return
		}
	}
}

var invocationCount = 0
var wg = NewMyWaitGroup()
var mutex sync.Mutex

func main() {
	wg.Add(200)
	for i := 0; i < 100; i++ {
		go fn()
	}
	wg.Wait()
	fmt.Printf("Invocation count = %d\n", invocationCount)
}

func fn() {
	mutex.Lock()
	invocationCount++
	wg.Done()
	mutex.Unlock()
}
