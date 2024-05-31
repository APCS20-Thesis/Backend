package gophish

import "time"

type T struct {
	Id               int       `json:"id"`
	Name             string    `json:"name"`
	InterfaceType    string    `json:"interface_type"`
	FromAddress      string    `json:"from_address"`
	Host             string    `json:"host"`
	Username         string    `json:"username"`
	Password         string    `json:"password"`
	IgnoreCertErrors bool      `json:"ignore_cert_errors"`
	ModifiedDate     time.Time `json:"modified_date"`
	Headers          []Header  `json:"headers"`
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
