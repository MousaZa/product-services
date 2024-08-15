package handlers

import (
	"errors"
	"github.com/MousaZa/product-services/product-api/data"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// swagger:route DELETE /products/{id} products deleteProduct
// responses:
//
//	201: noContent

// DeleteProduct deletes a product form the database
func (p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}
	p.l.Println("Handle DELETE Products", id)

	err = data.DeleteProduct(id)

	if errors.Is(err, data.ErrProductNotFound) {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
}
