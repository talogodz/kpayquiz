package main

import (
	"kpayquiz/database"
	"kpayquiz/merchant"

	"github.com/gin-gonic/gin"
)

func main() {
	dao := database.DAO{"localhost", "merchant_db"}
	dao.Connect()
	r := gin.Default()
	v1 := r.Group("/api/v1")
	v1.POST("/register", merchant.RegisterEndPoint)
	v1.POST("/buy/product", merchant.BuyProductEndPoint)

	merchant.RouteRegister(v1.Group("/merchant"))
	r.Run(":8080")
}
