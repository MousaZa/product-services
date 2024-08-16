package handlers

import (
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
	p.l.Info("Handle GET Products")
	cur := r.URL.Query().Get("currency")
	rw.Header().Add("Content-Type", "application/json")

	lp, _ := p.productDB.GetProducts(cur)
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
	cur := r.URL.Query().Get("currency")
	id, err := strconv.Atoi(vars["id"])
	rw.Header().Add("Content-Type", "application/json")
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}
	lp, _ := p.productDB.GetSingleProduct(id, cur)

	//p.l.Printf("Resp %#v", resp)
	//
	//lp.Price = lp.Price * resp.Rate

	err = lp.ToJSONSingle(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
