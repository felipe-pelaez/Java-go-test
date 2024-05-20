package main

import (
	handlerProduct "MicroserviceTemplate/cmd/handler/product"
	routerProduct "MicroserviceTemplate/cmd/router/product"
	"MicroserviceTemplate/config"
	_ "MicroserviceTemplate/docs"
	"MicroserviceTemplate/internal/product"
	"MicroserviceTemplate/pkg/eureka"
	"MicroserviceTemplate/pkg/middleware"
	store "MicroserviceTemplate/pkg/store/product"
	"context"
	_ "github.com/dimiro1/banner/autoload"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
	"log"
	"net"
	"net/http"
	"strconv"
)

// @title 							ms-template-mongo-go
// @version 						1.0.0
// @description 					This is a sample server for a microservice template in ecosystem Java Spring Cloud and Go.
// @contact.name 					Nelson David Camacho Ovalle
// @license.name 					Apache 2.0
// @license.url 					http://www.apache.org/licenses/LICENSE-2.0.html
// @host      		            	localhost:8080
// @securityDefinitions.apikey 		BearerAuth
// @in 								header
// @name 							Authorization
// @BasePath  						/
func main() {

	ctx := context.Background()

	_ = fx.New(
		fx.Provide(
			store.NewStore,
			product.NewRepository,
			product.NewService,
			handlerProduct.NewHandler,
			routerProduct.NewProductRouter,
		),
		fx.Invoke(
			LifecycleHooks,
		),
		fx.NopLogger,
	).Start(ctx)

}

// LifecycleHooks - Initializes application hooks in the application life cycle.
func LifecycleHooks(lc fx.Lifecycle, router routerProduct.IRouter) {
	lc.Append(fx.Hook{
		OnStart: func(c context.Context) error {

			// ? ================== Load configuration ================== ?

			vp := viper.New()

			vp.SetConfigName("application")
			vp.SetConfigType("yaml")
			vp.AddConfigPath("./resources")

			err := vp.ReadInConfig()
			if err != nil {
				log.Fatalln(err)
			}

			config.LoadConfigurationFromBranch(
				vp.GetString("application.config.import"),
				vp.GetString("application.name"),
				vp.GetString("application.config.profile"),
				vp.GetString("application.config.branch"),
			)

			// ==================== Start server ==================== ?

			appName := viper.GetString("application.name")
			port := viper.GetString("server.port")
			appId := uuid.New().String()

			gin.SetMode(gin.ReleaseMode)
			r := gin.Default()
			r.Use(middleware.IsAuthorizedJWT("/swagger/*any"))
			r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
			r = router.GetRoutes(r)

			ln, err := net.Listen("tcp", ":"+port)
			if err != nil {
				return err
			}

			_, portObtained, err := net.SplitHostPort(ln.Addr().String())
			if err != nil {
				return err
			}

			log.Printf("listening on port %s", portObtained)

			portObtainedInt, err := strconv.Atoi(portObtained)
			if err != nil {
				return err
			}

			eureka.StartClient(appName, appId, portObtainedInt)

			err = http.Serve(ln, r)
			if err != nil {
				return err
			}

			return nil

		},
		OnStop: func(c context.Context) error {
			log.Print("stopping...")
			return nil
		},
	})
}
