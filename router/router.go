package router

import (
	"net/http"
	"os"

	"go-ocr/boot"
	"go-ocr/infrastructure/config"
	"go-ocr/infrastructure/httplib"
	"go-ocr/infrastructure/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HandlerRouter struct {
	Setup boot.HandlerSetup
}

func NewHandlerRouter(setup boot.HandlerSetup) InterfaceRouter {
	return &HandlerRouter{
		Setup: setup,
	}
}

type InterfaceRouter interface {
	RouterWithMiddleware() *gin.Engine
}

func notFoundHandler(c *gin.Context) {
	// render 404 custom response
	httplib.SetErrorResponse(c, http.StatusNotFound, "Not Matching of Any Routes")
	return
}

func methodNotAllowedHandler(c *gin.Context) {
	// render 404 custom response
	httplib.SetErrorResponse(c, http.StatusMethodNotAllowed, "Method Not Allowed")
	return
}

func (hr *HandlerRouter) RouterWithMiddleware() *gin.Engine {
	//add new instance for bun router and add not found handler
	//and method with not allowed handler
	c := gin.New()

	//use recovery
	c.Use(gin.Recovery())

	//use options middleware handler
	c.OPTIONS("*any", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	//use cors need to updated if the requested need specific allow origin
	c.Use(cors.Default())

	//if logMode is true set logger to stdout on gin
	if config.Conf.LogMode {
		// Open a file to write logs to
		f, err := os.OpenFile("go-ocr.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}

		// Set the logger output to the file
		c.Use(gin.LoggerWithConfig(gin.LoggerConfig{
			Output: f,
		}))
	} else {
		// Log to stdout if LogMode is disabled
		c.Use(gin.LoggerWithConfig(gin.LoggerConfig{
			Output: gin.DefaultWriter,
		}))
	}

	//set middleware to use not found handler
	c.NoRoute(notFoundHandler)

	//set middleware to use method not allowed
	c.NoMethod(methodNotAllowedHandler)

	//serve static files
	c.Static("/uploads", "./uploads")

	//grouping on root endpoint
	api := c.Group("/api")

	api.Use(middleware.RateLimiterMiddleware(hr.Setup.Limiter))

	//grouping on "api/v1"
	v1 := api.Group("/v1")

	//module health
	prefixHealth := v1.Group("/health")
	hr.Setup.HealthHttp.GroupHealth(prefixHealth)

	//module ocr
	prefixOcr := v1.Group("/ocr")
	hr.Setup.OcrHttp.GroupOcr(prefixOcr)

	return c

}
