package data

import (
	"context"
	"encoding/json"
	"fmt"
	protos "github.com/MousaZa/product-services/currency/protos/currency"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-hclog"
	"io"
	"regexp"
	"time"
)

// Product defines the structure for an API product
// swagger:model
type Product struct {
	//	the id for the product
	//
	//	required: true
	// 	min: 1
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	err := e.Decode(p)
	if err != nil {
		// Log the error or handle it as needed

		return fmt.Errorf("error decoding JSON: %v", err)
	}
	return nil
}

func (p *Product) Validate() error {
	validate := validator.New()
	err := validate.RegisterValidation("sku", validateSKU)
	if err != nil {
		return err
	}
	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {

	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)
	if len(matches) != 1 {
		return false
	}
	return true
}

type Products []*Product

type ProductsDB struct {
	currency protos.CurrencyClient
	log      hclog.Logger
}

func NewProductsDB(client protos.CurrencyClient, l hclog.Logger) *ProductsDB {
	return &ProductsDB{log: l, currency: client}
}

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) ToJSONSingle(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *ProductsDB) GetProducts(currency string) (Products, error) {
	if currency == "" {
		return ProductList, nil
	}

	rate, err := p.GetRate(currency)
	if err != nil {
		p.log.Error("Unable to get rate", "currency", currency, "error", err)
		return nil, err
	}
	pr := Products{}
	for _, p := range ProductList {
		np := *p
		np.Price = np.Price * rate.Rate
		pr = append(pr, &np)
	}

	return pr, nil
}

func (p *ProductsDB) GetSingleProduct(id int, currency string) (*Product, error) {
	prod, _, err := FindProduct(id)
	if err != nil {
		return nil, err
	}
	if currency == "" {
		return prod, nil
	}
	rate, err := p.GetRate(currency)
	if err != nil {
		p.log.Error("Unable to get rate", "currency", currency, "error", err)
		return nil, err
	}
	np := *prod
	p.log.Info("rate", rate)
	fmt.Printf("rate %s", rate)
	np.Price = np.Price * rate.Rate
	return &np, nil
}

func DeleteProduct(id int) error {
	_, pos, err := FindProduct(id)
	if err != nil {
		return err
	}
	ProductList = append(ProductList[:pos], ProductList[pos+1:]...)
	return nil
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	ProductList = append(ProductList, p)
}
func getNextID() int {
	lp := ProductList[len(ProductList)-1]
	return lp.ID + 1
}
func (p *ProductsDB) UpdateProducts(id int, prod *Product) error {
	_, pos, err := FindProduct(id)
	if err != nil {
		return err
	}
	prod.ID = id
	ProductList[pos] = prod
	return nil
}

var ErrProductNotFound = fmt.Errorf("product not found")

func FindProduct(id int) (*Product, int, error) {
	for i, p := range ProductList {
		if p.ID == id {
			return p, i, nil
		}

	}
	return nil, 0, ErrProductNotFound
}

func (p *ProductsDB) GetRate(destination string) (*protos.RateResponse, error) {
	rr := &protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value[destination]),
	}

	resp, err := p.currency.GetRate(context.Background(), rr)
	p.log.Info("resp", resp)
	return resp, err
}

var ProductList = Products{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
