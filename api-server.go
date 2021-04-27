package main

import (
	"echoserver/echox"
	"encoding/json"
	"strconv"
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

func findProductAsync(rc echox.RoutingContext) error {
	strid := rc.Param("id")
	id, err := strconv.Atoi(strid)
	if err != nil {
		return rc.BadRequest()
	}
	result := <-repos.findByIdFuture(rc.Request().Context(), id)
	//rs, err := repos.findById(rc.Request().Context(), id)
	// result := <-r
	return rc.JsonWrap(result.Product, result.Error)
}

func findProduct(rc echox.RoutingContext) error {
	strid := rc.Param("id")
	id, err := strconv.Atoi(strid)
	if err != nil {
		return rc.BadRequest()
	}
	rs, err := repos.findById(rc.Request().Context(), id)
	return rc.JsonWrap(rs, err)
}

func createProduct(rc echox.RoutingContext) error {
	p := new(Product)
	if err := json.NewDecoder(rc.Request().Body).Decode(p); err != nil {
		rc.Logger().Error(err)
		return rc.BadRequest()
	}
	err := repos.create(rc.Request().Context(), p)
	return rc.JsonWrap(p, err)
}

func listProducts(rc echox.RoutingContext) error {
	rs, err := repos.list(rc.Request().Context())
	return rc.JsonWrap(rs, err)
}

func deleteProduct(rc echox.RoutingContext) error {
	strid := rc.Param("id")
	id, err := strconv.Atoi(strid)
	if err != nil {
		return rc.BadRequest()
	}
	err = repos.deleteById(rc.Request().Context(), id)
	return rc.JsonWrap(nil, err)
}

func updateProduct(rc echox.RoutingContext) error {
	id := rc.Param("id")
	iid, err := strconv.Atoi(id)
	if err != nil {
		return rc.BadRequest()
	}
	p := new(Product)
	if err := json.NewDecoder(rc.Request().Body).Decode(p); err != nil {
		rc.Logger().Error(err)
		return rc.BadRequest()
	}
	np, err := repos.updateById(rc.Request().Context(), iid, p)
	return rc.JsonWrap(np, err)
}
