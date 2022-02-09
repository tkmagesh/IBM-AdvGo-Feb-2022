package main

import (
	"encoding/json"
	"net/http"
)

type Product struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Cost  float32 `json:"cost"`
	Units int     `json:"units"`
}

var products = []Product{
	{Id: 1, Name: "Pen", Cost: 5.5, Units: 100},
	{Id: 2, Name: "Pencil", Cost: 2, Units: 50},
	{Id: 3, Name: "Marker", Cost: 12, Units: 10},
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	http.HandleFunc("/products", productsHandler)
	http.ListenAndServe(":8080", nil)
}

func productsHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		encoder := json.NewEncoder(rw)
		encoder.Encode(products)
	case "POST":
		decoder := json.NewDecoder(r.Body)
		var newProduct Product
		decoder.Decode(&newProduct)
		//assign a unique to the newProduct
		newProduct.Id = len(products) + 1
		products = append(products, newProduct)
		rw.WriteHeader(http.StatusCreated)
	case "PUT":
		rw.Write([]byte("The given product is updated"))
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}
