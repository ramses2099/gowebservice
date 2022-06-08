package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//REST (Representation State Transfer)

// Handling HTTP Requests
// -GET, PUT, POST, DELETE
// -Encoding Data

// HTTP Request
// - Method
// - Headers
// - Body

// HTTP Response
// - Status Code (200: Success)
// - Headers
// - Body

type Product struct {
	ProductID      int    `json:"productId"`
	Manufacturer   string `json:"manufacturer"`
	Sku            string `json:"sku"`
	Upc            string `json:"upc"`
	PricePerUnit   string `json:"pricePerUnit"`
	QuantityOnHand int    `json:"quantityOnHand"`
	ProductName    string `json:"productName"`
}

var productList []Product

func init() {
	productJSON := `[
		{
			"productId":1,
			"manufacturer":"Johns-Jenkins",
			"sku":"p5z3456ds",
			"upc":"95464446464646",
			"pricePerUnit":"457.135",
			"quantityOnHand":9703,
			"productName":"Lavanda"
		},
		{
			"productId":2,
			"manufacturer":"Johns-Jenkins",
			"sku":"p5z3456ds",
			"upc":"33345556",
			"pricePerUnit":"100.135",
			"quantityOnHand":73,
			"productName":"silas"
		},
		{
			"productId":3,
			"manufacturer":"Johns-Jenkins",
			"sku":"p5z345fsfs6ds",
			"upc":"8879546644646",
			"pricePerUnit":"75.00",
			"quantityOnHand":703,
			"productName":"Corales"
		}
	]`

	err := json.Unmarshal([]byte(productJSON), &productList)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	//http.HandleFunc
	http.HandleFunc("/products", productsHandler)
	http.HandleFunc("/product/", productHandler)

	// http.ListAndServeTLS for ssl
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}

}

//
func findProductById(productID int) (*Product, int) {
	for i, pr := range productList {
		if pr.ProductID == productID {
			return &pr, i
		}
	}
	return nil, 0
}

//
func getNextID() int {
	highestID := -1
	for _, product := range productList {
		if highestID < product.ProductID {
			highestID = product.ProductID
		}
	}
	return highestID + 1
}

//http.HandleFunc
func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productsJson, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("COntent-Type", "application/json")
		w.Write(productsJson)
	case http.MethodPost:
		// add a new product to the list
		var newProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//
		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//
		if newProduct.ProductID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return

		}
		//
		newProduct.ProductID = getNextID()
		productList = append(productList, newProduct)
		w.WriteHeader(http.StatusCreated)
		return
	}
}

//
func productHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "product/")
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//
	product, _ := findProductById(productID)
	//
	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		prodcutJson, err := json.Marshal(product)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("COntent-Type", "application/json")
		w.Write(prodcutJson)
	case http.MethodPut:

	}
}
