package main

import (
	"kpayquiz/merchant"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	merchant.RouteRegister(v1.Group("/merchant"))
	r.Run(":8080")
}
