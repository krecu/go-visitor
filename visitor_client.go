package visitor

import (
	"net/http"
	"time"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"viqeo-stats/visitor/model"
)

type Client struct {
	Addr string
}

type Body struct {
	Ip string
	Ua string
	Id string
}

func New(addr string) (*Client) {

	return &Client{
		Addr: addr,
	}
}

func (v *Client) Get(id string, ip string, ua string) (*model.Visitor, error){
	var proto model.Visitor

	url := v.Addr + "/api/visitor"

	body, err := json.Marshal(&Body{
		Ip: ip,
		Id: id,
		Ua: ua,
	})

	if err != nil {
		return &model.Visitor{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	var client = &http.Client{
		Timeout: time.Second * 100,
	}
	response, err := client.Do(req)

	if err != nil || response.StatusCode != 200 {
		return &model.Visitor{}, err
	}
	defer response.Body.Close()

	result, _ := ioutil.ReadAll(response.Body)

	err = json.Unmarshal(result, &proto)

	return &proto, err
}