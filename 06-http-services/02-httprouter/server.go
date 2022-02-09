package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
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
	router := httprouter.New()
	router.ServeFiles("/public/*filepath", http.Dir("./static"))
	router.GET("/", rootHandler)
	router.GET("/products", serveAllProducts)
	router.GET("/products/:id", serveProductById)
	router.POST("/products", addNewProduct)

	http.ListenAndServe(":8080", router)
}

func rootHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("Hello World!"))
}

func serveAllProducts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	encoder := json.NewEncoder(w)
	encoder.Encode(products)
}

func serveProductById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, _ := strconv.Atoi(params.ByName("id"))
	var product Product
	for _, p := range products {
		if p.Id == id {
			product = p
			break
		}
	}
	if product.Id == id {
		encoder := json.NewEncoder(w)
		encoder.Encode(product)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func addNewProduct(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var newProduct Product
	decoder.Decode(&newProduct)
	//assign a unique to the newProduct
	newProduct.Id = len(products) + 1
	products = append(products, newProduct)
	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.Encode(newProduct)

}
