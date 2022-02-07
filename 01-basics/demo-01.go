package main

import "fmt"

func main() {
	/* 1. Functions can be assigned to variables */
	/*
		var fn = func() {
			fmt.Println("fn is invoked")
		}
	*/
	/*
		fn := func() {
			fmt.Println("fn is invoked")
		}
	*/

	var fn func()
	fn = func() {
		fmt.Println("fn is invoked")
	}
	fn()

	/*
		add := func(x, y int) int {
			return x + y
		}
	*/
	var add func(int, int) int
	add = func(x, y int) int {
		return x + y
	}
	fmt.Println(add(100, 200))

	/* 2. Functions can be passed as arguments to other functions */
	exec(fn)
	execOper(add, 100, 200)

	/* 3. Functions can be retured from other function */
	adderFor100 := getAdderFor(100)
	fmt.Println(adderFor100(200))

	/* 4. Closures */
	increment := incrementor()
	fmt.Println(increment())
	fmt.Println(increment())
	fmt.Println(increment())
	fmt.Println(increment())

	/* Anonymous Function */
	func() {
		fmt.Println("Anonymous function invoked")
	}()
}

func exec(fn func()) {
	fn()
}

func execOper(operFn func(int, int) int, x, y int) {
	fmt.Println(operFn(x, y))
}

func getAdderFor(x int) func(int) int {
	return func(y int) int {
		return x + y
	}
}

/*
func getHttpClient(baseUrl string){
	return func(url, reqType, payload){

	}
}

httpClient(url, reqType, payload)

GET 	http://myNewService.com/products
GET 	http://myNewService.com/products/1
POST 	http://myNewService.com/products
PUT 	http://myNewService.com/products/1
DELETE 	http://myNewService.com/products/1
*/

/* Closures */
func incrementor() func() int {
	var count = 0
	return func() int {
		count++
		return count
	}
}
