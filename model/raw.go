package model


type Raw struct {
	Id string
	Created uint
	Updated uint
	City map[string]interface{}
	Country map[string]interface{}
	Location map[string]interface{}
	Postal map[string]interface{}
	Region map[string]interface{}
	Browser map[string]interface{}
	Device map[string]interface{}
	Platform map[string]interface{}
	Ip map[string]interface{}
	Extra map[string]interface{}
}