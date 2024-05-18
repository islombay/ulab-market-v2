package models

type Response struct {
	StatusCode  int `json:"status_code"`
	Description string
	Data        interface{} `json:"data"`
	Count       int         `json:"count"`
}
