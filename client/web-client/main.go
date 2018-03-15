package web_client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"fmt"

	"github.com/krecu/go-visitor/model"
)

type Client struct {
	Addr string
}

type Body struct {
	Ip string
	Ua string
	Id string
}

type Patch struct {
	Id     string
	Fields string
}

func New(addr string) *Client {

	return &Client{
		Addr: addr,
	}
}

func (v *Client) Get(id string) (proto *model.Raw, err error) {

	url := fmt.Sprintf("%s/api/visitor/%s", v.Addr, id)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")

	var client = &http.Client{
		Timeout: time.Second * 1,
	}
	response, err := client.Do(req)

	if err != nil || response.StatusCode != 200 {
		return &model.Raw{}, err
	}
	defer response.Body.Close()

	result, _ := ioutil.ReadAll(response.Body)

	err = json.Unmarshal(result, &proto)

	return
}

// создание на GRPC
func (v *Client) Post(id string, ip string, ua string, extra map[string]interface{}) (proto *model.Raw, err error) {

	url := fmt.Sprintf("%s/api/visitor", v.Addr)

	body, err := json.Marshal(&Body{
		Ip: ip,
		Id: id,
		Ua: ua,
	})

	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	var client = &http.Client{
		Timeout: time.Second * 1,
	}
	response, err := client.Do(req)

	if err != nil || response.StatusCode != 200 {
		return
	}
	defer response.Body.Close()

	result, _ := ioutil.ReadAll(response.Body)

	err = json.Unmarshal(result, &proto)

	return
}

// создание на GRPC
func (v *Client) Patch(id string, fields map[string]interface{}) (proto *model.Raw, err error) {

	var (
		filedJson []byte
	)

	filedJson, err = json.Marshal(fields)
	if err != nil {
		return
	}

	url := fmt.Sprintf("%s/api/visitor/%s", v.Addr, id)

	body, err := json.Marshal(&Patch{
		Id:     id,
		Fields: string(filedJson),
	})

	if err != nil {
		return
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	var client = &http.Client{
		Timeout: time.Second * 1,
	}
	response, err := client.Do(req)

	if err != nil || response.StatusCode != 200 {
		return
	}
	defer response.Body.Close()

	result, _ := ioutil.ReadAll(response.Body)

	err = json.Unmarshal(result, &proto)

	return
}
