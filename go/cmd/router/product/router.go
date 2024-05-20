package product

import (
	"MicroserviceTemplate/cmd/handler/product"
	"github.com/gin-gonic/gin"
)

// ? ==================== Interfaces ====================

type IRouter interface {
	GetRoutes(r *gin.Engine) *gin.Engine
}

// ? ==================== Structures ==================== ?

type Router struct {
	Handler product.IHandler
}

// ? ==================== Constructor ==================== ?

// NewProductRouter returns a new product router
func NewProductRouter(handler product.IHandler) IRouter {
	return &Router{handler}
}

// ? ===================== Methods ==================== ?

// GetRoutes returns product routes
func (router *Router) GetRoutes(r *gin.Engine) *gin.Engine {

	routerProducts := r.Group("/products")

	routerProducts.GET("/", router.Handler.GetAll())
	routerProducts.GET("/:id", router.Handler.GetByID())
	routerProducts.POST("/", router.Handler.Save())
	routerProducts.PUT("/:id", router.Handler.Update())
	routerProducts.PATCH("/:id", router.Handler.PatchUpdate())
	routerProducts.DELETE("/:id", router.Handler.Delete())

	return r

}
