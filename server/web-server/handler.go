package main

import (
	"github.com/valyala/fasthttp"
	"encoding/json"
	"time"
	"github.com/fatih/structs"
)

type Handler struct {}

func HewHandler() *Handler{
	return &Handler{}
}

func (h *Handler) GetHandler(ctx *fasthttp.RequestCtx) {
	h.CrossDomain(ctx)
	ctx.SetContentType("application/json; charset=utf8")
	ctx.SetStatusCode(200)
	ctx.Write([]byte("TEST"))
}

func (h *Handler) PostHandler(ctx *fasthttp.RequestCtx) {

	type BodyValues struct{
		Ip string
		Ua string
		Id string
		Extra map[string]interface{}
	}

	var body BodyValues

	json.Unmarshal(ctx.Request.Body(), &body)

	extra := make(map[string]interface{})
	extra["id"] = body.Id
	extra["created"] = time.Now().Unix()

	ctx.SetContentType("application/json; charset=utf8")

	info, err := Core.Get(body.Id, body.Ip, body.Ua, extra); if err == nil {

		m := structs.Map(info)
		for key, val := range extra {
			m[key] = val
		}

		jsonData, err := json.Marshal(m); if err == nil {
			ctx.SetStatusCode(200)
			ctx.Write(jsonData)
		}
	}

	if err != nil {
		ctx.SetStatusCode(504)
		ctx.Write([]byte("ERROR:" + err.Error()))
	}

}


func (h *Handler) PutHandler(ctx *fasthttp.RequestCtx) {
	h.CrossDomain(ctx)
	ctx.SetContentType("application/json; charset=utf8")
	ctx.SetStatusCode(200)
	ctx.Write([]byte("TEST"))
}


// добавляем кросориджены к респонсу
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