package product

import (
	"MicroserviceTemplate/internal/domain"
	"MicroserviceTemplate/internal/product"
	"MicroserviceTemplate/pkg/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ? ==================== Interfaces ====================

type IHandler interface {
	GetAll() gin.HandlerFunc
	GetByID() gin.HandlerFunc
	Save() gin.HandlerFunc
	Update() gin.HandlerFunc
	PatchUpdate() gin.HandlerFunc
	Delete() gin.HandlerFunc
}

// ? ==================== Structs ==================== ?

type Handler struct {
	service product.IService
}

// ? ==================== Constructors ==================== ?

// NewHandler retorna un nuevo handler de productos
func NewHandler(service product.IService) IHandler {
	return &Handler{service}
}

// ? ===================== Methods ==================== ?

// GetAll 		Returns all products
// @Summary 	Get all products
// @Tags 		Products
// @Description Get all products
// @Produce  	json
// @Success 	200 {object} domain.Products
// @Failure 	401 {object} web.ErrorResponse
// @Security    BearerAuth
// @Router 		/products [get]
func (handler *Handler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := handler.service.GetAll()
		if err != nil {
			web.ErrorResponseBody(c, http.StatusInternalServerError, "get_all_error", err.Error())
			return
		}
		web.SuccessResponseBody(c, http.StatusOK, products)
	}
}

// * =========== *

// GetByID 		Returns a product by its ID
// @Summary 	Get product by ID
// @Tags 		Products
// @Description Get product by ID
// @Param 		id path string true "Product ID"
// @Produce 	json
// @Success 	200 {object} domain.Product
// @Failure 	404 {object} web.ErrorResponse
// @Failure 	401 {object} web.ErrorResponse
// @Security 	BearerAuth
// @Router 		/products/{id} [get]
func (handler *Handler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		productById, err := handler.service.GetByID(id)
		if err != nil {
			web.ErrorResponseBody(c, http.StatusNotFound, "not_found", "Product not found")
			return
		}

		web.SuccessResponseBody(c, http.StatusOK, productById)

	}
}

// * =========== *

// Save 		saves a product
// @Summary 	Save a product
// @Tags 		Products
// @Description Save a product
// @Accept  	json
// @Param 		product body domain.Product true "Product to save"
// @Produce 	json
// @Success 	201 {object} domain.Product
// @Failure 	400 {object} web.ErrorResponse
// @Failure 	401 {object} web.ErrorResponse
// @Security 	BearerAuth
// @Router 		/products [post]
func (handler *Handler) Save() gin.HandlerFunc {
	return func(c *gin.Context) {

		var productToSave domain.Product

		if err := c.ShouldBindJSON(&productToSave); err != nil {
			web.ErrorResponseBody(c, http.StatusBadRequest, "invalid_json", err.Error())
			return
		}

		productSaved, err := handler.service.Save(&productToSave)
		if err != nil {
			web.ErrorResponseBody(c, http.StatusBadRequest, "save_error", err.Error())
			return
		}

		web.SuccessResponseBody(c, http.StatusCreated, productSaved)

	}
}

// * =========== *

// Update 		updates a product
// @Summary 	Update a product
// @Tags 		Products
// @Description Update a product
// @Param 		id path string true "Product ID"
// @Accept  	json
// @Param 		product body domain.Product true "Product to update"
// @Produce  	json
// @Success 	200 {object} domain.Product
// @Failure 	400 {object} web.ErrorResponse
// @Failure 	401 {object} web.ErrorResponse
// @Security 	BearerAuth
// @Router 		/products/{id} [put]
func (handler *Handler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		var productToUpdate domain.Product

		if err := c.ShouldBindJSON(&productToUpdate); err != nil {
			web.ErrorResponseBody(c, http.StatusBadRequest, "invalid_json", err.Error())
			return
		}

		productToUpdate.ID = id

		err := handler.service.Update(&productToUpdate)
		if err != nil {
			web.ErrorResponseBody(c, http.StatusBadRequest, "update_error", err.Error())
			return
		}

		web.SuccessResponseBody(c, http.StatusOK, "Product updated")

	}

}

// * =========== *

// PatchUpdate 	Partially update a product (only the fields that are sent)
// @Summary 	Patch update a product (only the fields that are sent)
// @Tags 		Products
// @Description Patch update a product (only the fields that are sent)
// @Accept  	json
// @Param 		product body domain.Product true "Product to patch"
// @Param 		id path string true "Product ID"
// @Produce  	json
// @Success 	200 {object} domain.Product
// @Failure 	400 {object} web.ErrorResponse
// @Failure 	401 {object} web.ErrorResponse
// @Security 	BearerAuth
// @Router 		/products/{id} [patch]
func (handler *Handler) PatchUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		var productToUpdate domain.Product

		if err := c.ShouldBindJSON(&productToUpdate); err != nil {
			web.ErrorResponseBody(c, http.StatusBadRequest, "invalid_json", err.Error())
			return
		}

		productToUpdate.ID = id

		err := handler.service.PatchUpdate(&productToUpdate)
		if err != nil {
			web.ErrorResponseBody(c, http.StatusBadRequest, "update_error", err.Error())
			return
		}

		web.SuccessResponseBody(c, http.StatusOK, "Product updated")

	}

}

// * =========== *

// Delete 		deletes a product
// @Summary 	Delete a product
// @Tags 		Products
// @Description Delete a product
// @Param 		id path string true "Product ID"
// @Produce 	json
// @Success 	204 {object} domain.Product
// @Failure 	400 {object} web.ErrorResponse
// @Failure 	401 {object} web.ErrorResponse
// @Security 	BearerAuth
// @Router 		/products/{id} [delete]
func (handler *Handler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		err := handler.service.Delete(id)
		if err != nil {
			web.ErrorResponseBody(c, http.StatusBadRequest, "delete_error", err.Error())
			return
		}

		web.SuccessResponseBody(c, http.StatusOK, "Product deleted")

	}
}
