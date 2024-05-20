package product

import (
	"MicroserviceTemplate/internal/domain"
	"MicroserviceTemplate/internal/product"
	store "MicroserviceTemplate/pkg/store/product"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Product Service Suite")
}

var (
	productTest = domain.Product{
		ID:       uuid.New().String(),
		Name:     "Test",
		Price:    1.0,
		Quantity: 1,
	}
)

var _ = Describe("Product Service", func() {

	productStore := store.NewStore()
	productRepository := product.NewRepository(productStore)
	productService := product.NewService(productRepository)

	It("Save product", func() {

		product, err := productService.Save(&productTest)

		Expect(err).To(BeNil())
		Expect(err).NotTo(HaveOccurred())

		Expect(product.Name).To(Equal(productTest.Name))
		Expect(product.Price).To(Equal(productTest.Price))
		Expect(product.Quantity).To(Equal(productTest.Quantity))

	})

	It("Read all products", func() {

		products, err := productService.GetAll()

		Expect(err).To(BeNil())
		Expect(err).NotTo(HaveOccurred())

		Expect(len(*products) > 0).To(BeTrue())

		Expect((*products)[0].Name).To(Equal(productTest.Name))
		Expect((*products)[0].Price).To(Equal(productTest.Price))
		Expect((*products)[0].Quantity).To(Equal(productTest.Quantity))

	})

	It("Read product by id", func() {

		product, err := productService.GetByID(productTest.ID)

		Expect(err).To(BeNil())
		Expect(err).NotTo(HaveOccurred())

		Expect(product.Name).To(Equal(productTest.Name))
		Expect(product.Price).To(Equal(productTest.Price))
		Expect(product.Quantity).To(Equal(productTest.Quantity))

	})

	It("Update product", func() {

		productTest.Name = "Test Update"
		productTest.Price = 2.0
		productTest.Quantity = 2

		err := productService.Update(&productTest)

		Expect(err).To(BeNil())
		Expect(err).NotTo(HaveOccurred())

		product, err := productService.GetByID(productTest.ID)

		Expect(err).To(BeNil())
		Expect(err).NotTo(HaveOccurred())

		Expect(product.Name).To(Equal(productTest.Name))
		Expect(product.Price).To(Equal(productTest.Price))
		Expect(product.Quantity).To(Equal(productTest.Quantity))

	})

	It("Delete product", func() {

		err := productService.Delete(productTest.ID)

		Expect(err).To(BeNil())
		Expect(err).NotTo(HaveOccurred())

		_, err = productService.GetByID(productTest.ID)

		Expect(err).To(HaveOccurred())

	})

})
