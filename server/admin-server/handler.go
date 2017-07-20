package main

import (
	"github.com/valyala/fasthttp"
	"html/template"
)

type Handler struct {}

func HewHandler() *Handler{
	return &Handler{}
}

/**
	Главная страница админки
 */
func (h *Handler) AdminPage(ctx *fasthttp.RequestCtx){

	ctx.SetContentType("text/html; charset=utf8")

	t, _ := template.ParseFiles(Conf.Path + "/template/index.html")
	t.Execute(ctx, []byte("base"))
}

func (h *Handler) HomePage(ctx *fasthttp.RequestCtx){
	ctx.SetContentType("application/json; charset=utf8")
}

func (h *Handler) GetList(ctx *fasthttp.RequestCtx){
	ctx.SetContentType("application/json; charset=utf8")


}


/**
	добавляем кросориджены к респонсу
 */
func (h *Handler) CrossDomain(ctx *fasthttp.RequestCtx) {
	ref := ctx.Request.Header.Peek("Origin")
	if ref == nil {
		ref = ctx.Request.Header.Peek("Referer")
		if ref == nil {
			ref = []byte("*")
		}
	}


	ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
	ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
	ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
	ctx.Response.Header.Set("Access-Control-Allow-Origin", string(ref))
}