package product

import "MicroserviceTemplate/internal/domain"

// ? ====================== Interfaces ====================== ?

type IService interface {
	GetAll() (*domain.Products, error)
	GetByID(id string) (*domain.Product, error)
	Save(product *domain.Product) (domain.Product, error)
	Update(product *domain.Product) error
	PatchUpdate(product *domain.Product) error
	Delete(id string) error
}

// ? ====================== Estructuras ====================== ?

type Service struct {
	repository IRepository
}

// ? ====================== Structs ====================== ?

func NewService(repository IRepository) IService {
	return &Service{repository}
}

// ? ====================== Methods ====================== ?

// GetAll returns all products
func (s *Service) GetAll() (*domain.Products, error) {
	return s.repository.GetAll()
}

// * =========== *

// GetByID returns a product by its ID
func (s *Service) GetByID(id string) (*domain.Product, error) {
	return s.repository.GetByID(id)
}

// * =========== *

// Save saves a product
func (s *Service) Save(product *domain.Product) (domain.Product, error) {
	return s.repository.Save(product)
}

// * =========== *

// Update update a product
func (s *Service) Update(product *domain.Product) error {
	return s.repository.Update(product)
}

// * =========== *

// PatchUpdate partially update a product (only the fields to be sent)
func (s *Service) PatchUpdate(product *domain.Product) error {
	return s.repository.PatchUpdate(product)
}

// * =========== *

// Delete eliminates a product by its ID
func (s *Service) Delete(id string) error {
	return s.repository.Delete(id)
}
