package domain

// ? =================== Structs =================== ?

type Product struct {
	ID       string  `bson:"_id" json:"_id"`
	Name     string  `bson:"name" json:"name"`
	Quantity int     `bson:"quantity" json:"quantity"`
	Price    float64 `bson:"price" json:"price"`
}

// * =========== *

type Products []*Product

// ? =================== Constructors =================== ?

func NewProduct(id string, name string, quantity int, price float64) *Product {
	return &Product{
		ID:       id,
		Name:     name,
		Quantity: quantity,
		Price:    price,
	}
}

// * =========== *

func NewProducts() *Products {
	return &Products{}
}
