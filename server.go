package main

import (
	"echoserver/echox"
)

func startAPIServer() {
	ex := echox.New()

	ex.POST("/products", createProduct)
	ex.GET("/products", listProducts)
	ex.GET("/products/:id", findProduct)
	ex.GET("/async/products/:id", findProductAsync)
	ex.PUT("/products/:id", updateProduct)
	ex.DELETE("/products/:id", deleteProduct)

	ex.Serve()
}
