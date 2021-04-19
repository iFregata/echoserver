package main

import (
	"echoserver/echolet"
	"encoding/json"
	"strconv"
)

var repos Repos

func init() {
	repos = CreateRepos()
}

func startAPIServer() {
	el := echolet.New()

	el.POST("/products", createProduct)
	el.GET("/products", listProducts)
	el.PUT("/products/:id", updateProduct)
	el.DELETE("/products/:id", deleteProduct)

	el.Serve()
}

func createProduct(rc echolet.RoutingContext) error {
	p := new(Product)
	if err := json.NewDecoder(rc.Request().Body).Decode(p); err != nil {
		rc.Logger().Error(err)
		return rc.BadRequest()
	}
	err := repos.create(rc.Request().Context(), p)
	return rc.RespJson(p, err)
}

func listProducts(rc echolet.RoutingContext) error {
	rs, err := repos.list(rc.Request().Context())
	return rc.RespJson(rs, err)
}

func deleteProduct(rc echolet.RoutingContext) error {
	strid := rc.Param("id")
	id, err := strconv.Atoi(strid)
	if err != nil {
		return rc.BadRequest()
	}
	err = repos.deleteById(rc.Request().Context(), id)
	return rc.RespJson(nil, err)
}

func updateProduct(rc echolet.RoutingContext) error {
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
	return rc.RespJson(np, err)
}
