package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

const (
	corsAllowHeaders     = "authorization"
	corsAllowMethods     ="HEAD,GET,POST,PUT,DELETE,OPTIONS"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "true"
)


type Http struct {
}

func NewHttp(addr string) (error){

	router := fasthttprouter.New()
	handler := HewHandler()

	// учет статистики
	router.GET("/api/visitor", handler.GetHandler)
	router.POST("/api/visitor", handler.PostHandler)
	router.PUT("/api/visitor", handler.PutHandler)

	return fasthttp.ListenAndServe(addr, router.Handler)
}