package handlers

import (
	"errors"
	"github.com/MousaZa/product-services/product-api/data"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// swagger:route PUT /products/{id} products updateProduct
// responses:
//
//	201: noContent

// UpdateProducts updates a product in the database
func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}
	p.l.Debug("Handle PUT Products", id)
	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	err = p.productDB.UpdateProducts(id, prod)
	if errors.Is(err, data.ErrProductNotFound) {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
}
