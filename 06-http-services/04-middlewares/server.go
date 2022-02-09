package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func main() {
	/*
		http.HandleFunc("/foo", appInt(profile(logger(handleFoo))))
		http.HandleFunc("/bar", appInt(profile(logger(handleBar))))
	*/
	http.HandleFunc("/foo", chain(handleFoo, logger, profile, appInt))
	http.HandleFunc("/bar", chain(handleBar, logger, profile, appInt))
	http.ListenAndServe(":8080", nil)
}

/* func logger(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s\t%s\n", r.Method, r.URL)
		handler(w, r)
	}
} */

func chain(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		handler = m(handler)
	}
	return handler
}

func profile(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			end := time.Now()
			elapsed := end.Sub(start)
			fmt.Printf("time taken : %v\n", elapsed)
		}()
		handler(w, r)
	}
}

func appInt(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctxWithValue := context.WithValue(ctx, "app-id", "asdfjklsfd")
		req := r.WithContext(ctxWithValue)
		handler(w, req)
	}
}

func logger(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s\t%s\n", r.Method, r.URL)
		handler(w, r)
	}
}

func handleFoo(w http.ResponseWriter, r *http.Request) {
	time.Sleep(500 * time.Millisecond)
	fmt.Println(r.Context().Value("app-id"))
	w.Write([]byte("foo"))
}

func handleBar(w http.ResponseWriter, r *http.Request) {
	time.Sleep(300 * time.Millisecond)
	fmt.Println(r.Context().Value("app-id"))
	w.Write([]byte("bar"))
}
