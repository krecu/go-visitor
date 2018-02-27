package model

type System struct {
	Device   Device   `json:"device"`
	Browser  Browser  `json:"browser"`
	Platform Platform `json:"platform"`
}
