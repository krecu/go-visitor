package main

type Record struct {
	Id    string                 `json:"id"`
	Ip    string                 `json:"ip"`
	Ua    string                 `json:"ua"`
	Extra map[string]interface{} `json:"extra"`
}
