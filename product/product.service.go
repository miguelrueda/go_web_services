package product

import (
	"encoding/json"
	"fmt"
	"go_web_services/cors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/websocket"
)

const productsBasePath = "products"

// SetupRoutes : Set the routes to add the handlers
func SetupRoutes(apiBasePath string) {
	handleProducts := http.HandlerFunc(productsHandler)
	handleProduct := http.HandlerFunc(productHandler)
	reportHandler := http.HandlerFunc(handleProductReport)
	http.Handle("/websocket", websocket.Handler(productSocket))
	fmt.Printf("Setting up %s/%s\n", apiBasePath, productsBasePath)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, productsBasePath), cors.Middleware(handleProducts))
	fmt.Printf("Setting up %s/%s/\n", apiBasePath, productsBasePath)
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, productsBasePath), cors.Middleware(handleProduct))
	fmt.Printf("Setting up %s/%s/reports\n", apiBasePath, productsBasePath)
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, productsBasePath), cors.Middleware(reportHandler))
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "products/")
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	product, err := getProduct(productID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:

		// return a single product
		productJSON, err := json.Marshal(product)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(productJSON)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPut:
		var product Product
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if product.ProductID != productID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// product = &product
		err = updateProduct(product)
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		removeProduct(productID)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func productsHandler(e http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productList, err := getProductList()
		if err != nil {
			e.WriteHeader(http.StatusInternalServerError)
		}
		productsJSON, err := json.Marshal(productList)
		if err != nil {
			e.WriteHeader(http.StatusInternalServerError)
		}
		e.Header().Set("Content-Type", "application/json")
		e.Write(productsJSON)
	case http.MethodPost:
		var newProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			e.WriteHeader(http.StatusBadRequest)
		}
		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil {
			e.WriteHeader(http.StatusBadRequest)
			return
		}
		if newProduct.ProductID != 0 {
			e.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = insertProduct(newProduct)
		if err != nil {
			e.WriteHeader(http.StatusInternalServerError)
			return
		}
		e.WriteHeader(http.StatusCreated)
		return
	case http.MethodOptions:
		return
	}
}
