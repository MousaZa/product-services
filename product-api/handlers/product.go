// Package classification of Product API
//
// Documentation for Product API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	protos "github.com/MousaZa/product-services/currency/protos/currency"
	"github.com/MousaZa/product-services/product-api/data"
	"log"
	"net/http"
)

// A list of products returns in the response
// swagger:response productsResponse
type productsResponse struct {
	// All products in the system
	// in: body
	Body []data.Product
}

// A product returns in the response
// swagger:response productResponse
type productResponse struct {
	// All products in the system
	// in: body
	Body data.Product
}

// swagger:response noContent
type productsNoContent struct {
}

// swagger:parameters deleteProduct updateProduct listSingleProduct
type productIdParameterWrapper struct {
	// the id of the product to delete from the database
	// in: path
	// required: true
	ID int `json:"id"`
}

// swagger:parameters addProduct updateProduct
type productParameterWrapper struct {
	// The product to add to the database
	// in: body
	// required: true
	Product data.Product `json:"product"`
}

type Products struct {
	l  *log.Logger
	cc protos.CurrencyClient
}

func NewProducts(l *log.Logger, cc protos.CurrencyClient) *Products {
	return &Products{l, cc}
}

type KeyProduct struct{}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		// validate the product

		err = prod.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating product", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating product: %s", err),
				http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}
