package model

type ShippingAddress struct {
	Address
	Fullname        string `json:"fullname,omitempty"`
	IsResidential   bool   `json:"isResidential,omitempty"`
	RequestLiftgate bool   `json:"requestLiftgate,omitempty"`
}
