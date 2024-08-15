package handlers

import (
	"github.com/MousaZa/product-services/product-api/data"
	"net/http"
)

// swagger:route POST /products products addProduct
// responses:
//
//	201: noContent

// AddProduct adds a product to the database00
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products")
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(prod)
}
