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

	router.ServeFiles("/assets/*filepath", Conf.Path + "/template/assets")

	router.GET("/admin", handler.AdminPage)
	router.GET("/", handler.HomePage)

	//router.GET("/api/debug", handler.ApiDebug)
	//router.GET("/api/commits", handler.ApiGetCommits)
	//router.GET("/api/tags", handler.ApiGetTags)
	//router.GET("/api/branch", handler.ApiGetBranches)
	//
	//router.POST("/api/custom/save", handler.ApiSaveCustom)
	//router.PATCH("/api/custom/:name", handler.ApiUpdateCustom)
	//router.GET("/api/custom/list", handler.ApiListCustom)
	//router.POST("/api/custom/remove", handler.ApiRemoveCustom)
	//router.POST("/api/custom/build", handler.ApiBuildCustom)
	//
	//router.POST("/api/checkout/:hash", handler.ApiCheckout)
	//
	//router.GET("/js/:module", handler.ApiGetModule)

	return fasthttp.ListenAndServe(addr, router.Handler)
}