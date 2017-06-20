package main

import (
	"github.com/valyala/fasthttp"
	"encoding/json"
)

type Handler struct {}

func HewHandler() *Handler{
	return &Handler{}
}

/**
	Пытаемся получить данные о визиотре
 */
func (h *Handler) GetHandler(ctx *fasthttp.RequestCtx) {

	type BodyValues struct{
		Ip string
		Ua string
		Id string
		Extra map[string]interface{}
	}

	var body BodyValues
	var err error

	json.Unmarshal(ctx.Request.Body(), &body)

	ctx.SetContentType("application/json; charset=utf8")

	if body.Id == "" || body.Ua == "" || body.Ip == "" {
		ctx.SetStatusCode(400)
		ctx.Write([]byte("ERROR: Bad parameters"))
		return
	} else {

		info, err := Core.Get(body.Id, body.Ip, body.Ua, body.Extra); if err == nil {

			jsonData, err := json.Marshal(info); if err == nil {
				ctx.SetStatusCode(200)
				ctx.Write(jsonData)
			}
		}
	}

	if err != nil {
		ctx.SetStatusCode(500)
		ctx.Write([]byte("ERROR:" + err.Error()))
	}
}

/**
	Пытаемся дополнить визитора наборами полей
 */
func (h *Handler) PutHandler(ctx *fasthttp.RequestCtx) {
	type BodyValues struct{
		Ip string
		Ua string
		Id string
		Extra map[string]interface{}
	}

	var body BodyValues
	var err error

	json.Unmarshal(ctx.Request.Body(), &body)

	ctx.SetContentType("application/json; charset=utf8")

	if body.Id == "" || body.Ua == "" || body.Ip == "" {
		ctx.SetStatusCode(400)
		ctx.Write([]byte("ERROR: Bad parameters"))
		return
	} else {

		info, err := Core.Put(body.Id, body.Ip, body.Ua, body.Extra); if err == nil {

			jsonData, err := json.Marshal(info); if err == nil {
				ctx.SetStatusCode(200)
				ctx.Write(jsonData)
			}
		}
	}

	if err != nil {
		ctx.SetStatusCode(500)
		ctx.Write([]byte("ERROR:" + err.Error()))
	}
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