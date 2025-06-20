package model

import "fmt"

var (
	ErrAddressNotFound  = fmt.Errorf("address not found")
	ErrDuplicateAddress = fmt.Errorf("address already exists")
	ErrInvalidAddress   = fmt.Errorf("address is invalid")
)

type Address struct {
	Id              int     `json:"id,omitempty" db:"id"`
	CustomerId      int     `json:"customerId,omitempty" db:"customer_id"`
	PostalCode      *string `json:"postalCode,omitempty" db:"postal_code"`
	CountryCode     string  `json:"countryCode,omitempty" db:"country_code"`
	SubdivisionCode *string `json:"subdivisionCode,omitempty" db:"subdivision_code"`
	SubdivisionName *string `json:"subdivisionName,omitempty" db:"subdivision_name"`
	CityName        string  `json:"cityName,omitempty" db:"city_name"`
	AddressLine1    string  `json:"addressLine1,omitempty" db:"address_line1"`
	AddressLine2    *string `json:"addressLine2,omitempty" db:"address_line2"`
}
