package merchant

import (
	"github.com/gin-gonic/gin"
)

func RouteRegister(r *gin.RouterGroup) {
	r.Use(BasicAuth)
	r.GET("/:id", InfoEndPoint)
	r.POST("/:id", UpdateEndPoint)
	r.GET("/:id/products", ProductsEndPoint)
	r.POST("/:id/product", AddProductEndPoint)
	r.POST("/:id/product/:product_id", UpdateProductEndPoint)
	r.DELETE("/:id/product/:product_id", RemoveProductEndPoint)
	r.POST("/:id/report", SellReportEndPoint)

}
