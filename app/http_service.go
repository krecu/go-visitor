package main

import (
	"encoding/json"
	"net/http"
	"time"

	"io/ioutil"

	"net"

	"github.com/CossackPyra/pyraconv"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

type HttpService struct {
	server *http.Server
	app    *App
}

//	Инициализурем вебсервер
func NewHttpService(app *App) (proto *HttpService, err error) {

	proto = &HttpService{
		app: app,
	}

	// устанавливаем корсы
	CorsMiddleWare := cors.New(cors.Options{
		AllowedOrigins:     proto.app.config.GetStringSlice("app.server.http.cors.AllowedOrigins"),
		AllowCredentials:   proto.app.config.GetBool("app.server.http.cors.AllowCredentials"),
		AllowedMethods:     proto.app.config.GetStringSlice("app.server.http.cors.AllowedMethods"),
		AllowedHeaders:     proto.app.config.GetStringSlice("app.server.http.cors.AllowedHeaders"),
		MaxAge:             proto.app.config.GetInt("app.server.http.cors.MaxAge"),
		Debug:              proto.app.config.GetBool("app.server.http.cors.Debug"),
		OptionsPassthrough: proto.app.config.GetBool("app.server.http.cors.OptionsPassthrough"),
	})

	// инициализируем роуты
	route := mux.NewRouter()

	route.HandleFunc("/{id}", proto.Response(proto.Get)).Methods("GET")
	route.HandleFunc("/", proto.Response(proto.Post)).Methods("POST")
	route.HandleFunc("/{id}", proto.Response(proto.Delete)).Methods("DELETE")
	route.HandleFunc("/{id}", proto.Response(proto.Patch)).Methods("PATCH")

	proto.server = &http.Server{
		ReadTimeout:    proto.app.config.GetDuration("app.server.http.ReadTimeout") * time.Second,
		WriteTimeout:   proto.app.config.GetDuration("app.server.http.WriteTimeout") * time.Second,
		MaxHeaderBytes: proto.app.config.GetInt("app.server.http.MaxHeaderBytes"),
		Addr:           proto.app.config.GetString("app.server.http.listen"),
		Handler:        CorsMiddleWare.Handler(route),
	}

	return
}

// старт вебсервера
func (h *HttpService) Start() {
	go func() {
		Logger.Fatal(h.server.ListenAndServe())
	}()
}

// получение модели
func (h *HttpService) Get(w http.ResponseWriter, r *http.Request) (response interface{}, code string) {

	_total := time.Now()

	vars := mux.Vars(r)

	if values, err := h.app.visitor.Get(pyraconv.ToString(vars["id"])); err == nil {
		response = values
		if values == nil {
			code = "0001"
		} else {
			code = "0000"
		}
	} else {

		if err == VisitorErrorEmpty {
			code = "0001"
		} else {
			code = err.Error()
		}
	}

	Logger.WithFields(logrus.Fields{
		"op":       "Get",
		"duration": time.Since(_total).Seconds(),
	}).Debugf("Get: %f", time.Since(_total).Seconds())

	return
}

func (h *HttpService) Post(w http.ResponseWriter, r *http.Request) (response interface{}, code string) {

	_total := time.Now()

	buf, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		code = err.Error()
		return
	}

	var msg struct {
		Id    string
		Ip    string
		Ua    string
		Extra map[string]interface{}
	}
	err = json.Unmarshal(buf, &msg)
	if err != nil {
		code = err.Error()
		return
	}

	if msg.Id == "" || msg.Ua == "" || msg.Ip == "" || net.ParseIP(msg.Ip) == nil {
		code = "0002"
		return
	}

	if values, err := h.app.visitor.Post(msg.Id, msg.Ip, msg.Ua, msg.Extra); err == nil {
		response = values
		code = "0000"
	} else {

		if err == VisitorErrorEmpty {
			code = "0001"
		} else {
			code = err.Error()
		}
	}

	Logger.WithFields(logrus.Fields{
		"op":       "Indent",
		"duration": time.Since(_total).Seconds(),
	}).Debugf("Indent: %f", time.Since(_total).Seconds())

	return
}

func (h *HttpService) Patch(w http.ResponseWriter, r *http.Request) (response interface{}, code string) {

	_total := time.Now()

	vars := mux.Vars(r)
	fields := make(map[string]interface{})

	buf, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		code = err.Error()
		return
	}

	err = json.Unmarshal(buf, &fields)
	if err != nil {
		code = err.Error()
		return
	}

	if fields == nil {
		code = "0002"
		return
	}

	if values, err := h.app.visitor.Patch(vars["id"], fields); err == nil {
		response = values
		code = "0000"
	} else {

		if err == VisitorErrorEmpty {
			code = "0001"
		} else {
			code = err.Error()
		}
	}

	Logger.WithFields(logrus.Fields{
		"op":       "Patch",
		"duration": time.Since(_total).Seconds(),
	}).Debugf("Patch: %f", time.Since(_total).Seconds())

	return
}

// удаление модели
func (h *HttpService) Delete(w http.ResponseWriter, r *http.Request) (response interface{}, code string) {

	_total := time.Now()

	vars := mux.Vars(r)

	if err := h.app.visitor.Delete(pyraconv.ToString(vars["id"])); err != nil {
		code = err.Error()
	} else {
		code = "0000"
	}

	Logger.WithFields(logrus.Fields{
		"op":       "Delete",
		"duration": time.Since(_total).Seconds(),
	}).Debugf("Delete: %f", time.Since(_total).Seconds())

	return
}

// хендлер для ловли ошибок
func (h *HttpService) Response(f func(http.ResponseWriter, *http.Request) (result interface{}, code string)) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		result, code := f(w, r)
		w.Header().Set("Content-Type", "application/json")

		switch code {
		case "0000":
			w.WriteHeader(http.StatusOK)
			break
		case "0001":
			w.WriteHeader(http.StatusNoContent)
			break
		case "0002":
			w.WriteHeader(http.StatusBadRequest)
			break
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		response := make(map[string]interface{})
		response["error"] = code
		response["result"] = result

		jsonData, _ := json.Marshal(response)

		w.Write(jsonData)

	})
}
