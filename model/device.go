package model

type Device struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Brand   string `json:"brand"`
	Mapping int    `json:"mapping_id"`
}
