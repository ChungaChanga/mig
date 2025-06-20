package model

type BillingAddress struct {
	Address
	Fullname string `json:"fullname,omitempty"`
}
