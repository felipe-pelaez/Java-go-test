package product

import (
	"MicroserviceTemplate/internal/domain"
	store "MicroserviceTemplate/pkg/store/product"
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// ? ==================== Interfaces ==================== ?

type IRepository interface {
	GetAll() (*domain.Products, error)
	GetByID(id string) (*domain.Product, error)
	Save(product *domain.Product) (domain.Product, error)
	Update(product *domain.Product) error
	PatchUpdate(product *domain.Product) error
	Delete(id string) error
}

// ? ==================== Structs ======================== ?

type Repository struct {
	db  *mongo.Collection
	ctx context.Context
}

// ? ==================== Constructors ==================== ?

// NewRepository returns a new product repository
func NewRepository(store store.IProductStore) IRepository {

	db, err := store.InitDatabase("products")
	if err != nil {
		log.Fatal(err)
	}

	return &Repository{db, context.Background()}
}

// ? ==================== Methods ====================== ?

// GetAll returns all products
func (r *Repository) GetAll() (*domain.Products, error) {

	var products domain.Products
	filter := bson.D{}

	cur, err := r.db.Find(r.ctx, filter)
	if err != nil {
		return nil, err
	}

	for cur.Next(r.ctx) {
		var product domain.Product
		err := cur.Decode(&product)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	err = cur.Close(r.ctx)
	if err != nil {
		return nil, err
	}

	return &products, nil

}

// * =========== *

// GetByID returns a product by its ID
func (r *Repository) GetByID(id string) (*domain.Product, error) {

	var product domain.Product
	filter := bson.D{{Key: "_id", Value: id}}

	err := r.db.FindOne(r.ctx, filter).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil

}

// * =========== *

// Save saves a product
func (r *Repository) Save(product *domain.Product) (domain.Product, error) {

	product.ID = uuid.New().String()

	_, err := r.db.InsertOne(r.ctx, product)
	if err != nil {
		return domain.Product{}, err
	}

	return *product, nil

}

// * =========== *

// Update update a product
func (r *Repository) Update(product *domain.Product) error {

	productToUpdate, err := r.GetByID(product.ID)
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", productToUpdate.ID}}

	update := bson.M{
		"$set": bson.M{
			"name":     product.Name,
			"quantity": product.Quantity,
		},
	}

	_, err = r.db.UpdateOne(r.ctx, filter, update)
	if err != nil {
		return err
	}

	return nil

}

// * =========== *

// PatchUpdate update a product partially with the fields that are sent to you
func (r *Repository) PatchUpdate(product *domain.Product) error {

	productToUpdate, err := r.GetByID(product.ID)
	if err != nil {
		return err
	}

	if product.Name != "" {
		productToUpdate.Name = product.Name
	}

	if product.Price != 0 {
		productToUpdate.Price = product.Price
	}

	if product.Quantity != 0 {
		productToUpdate.Quantity = product.Quantity
	}

	return r.Update(productToUpdate)

}

// * =========== *

// Delete eliminates a product
func (r *Repository) Delete(id string) error {

	productToDelete, err := r.GetByID(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": productToDelete.ID}

	_, err = r.db.DeleteOne(r.ctx, filter)
	if err != nil {
		return err
	}

	return nil

}
