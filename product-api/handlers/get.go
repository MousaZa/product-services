package handlers

import (
	"github.com/MousaZa/product-services/product-api/data"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
//
//	200: productsResponse
func (p *Products) GetAllProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	rw.Header().Add("Content-Type", "application/json")

	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// swagger:route GET /products/{id} products listSingleProduct
// Returns a product
// responses:
//
//	200: productResponse
func (p *Products) GetSingleProcut(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	rw.Header().Add("Content-Type", "application/json")
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}
	lp, _ := data.GetSingleProduct(id)
	err = lp.ToJSONSingle(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
