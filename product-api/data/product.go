package data

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
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

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) ToJSONSingle(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func GetProducts() Products {
	return ProductList
}

func GetSingleProduct(id int) (*Product, error) {
	prod, _, err := FindProduct(id)
	if err != nil {
		return nil, err
	}
	return prod, nil
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
func UpdateProducts(id int, p *Product) error {
	_, pos, err := FindProduct(id)
	if err != nil {
		return err
	}
	p.ID = id
	ProductList[pos] = p
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
