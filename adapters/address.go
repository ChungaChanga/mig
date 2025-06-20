package adapters

import (
	"context"

	"mig/address/model"
	addressProto "mig/api/ausweis/proto/address"
)

func NewAddressService(
	client addressProto.AddressServiceClient,
) *addressService {
	return &addressService{client: client}
}

type addressService struct {
	client addressProto.AddressServiceClient
}

func (s *addressService) CreateBillingAddress(ctx context.Context, customerId int, address *model.BillingAddress) error {
	_, err := s.client.CreateBillingAddress(ctx, &addressProto.BillingAddressRequest{
		CustomerId: int64(customerId),
		Address:    toBillingProto(address),
	})
	return err
}

func (s *addressService) CreateShippingAddress(ctx context.Context, customerId int, address *model.ShippingAddress) error {
	_, err := s.client.CreateShippingAddress(ctx, &addressProto.ShippingAddressRequest{
		CustomerId: int64(customerId),
		Address:    toShippingProto(address),
	})
	return err
}

func toShippingProto(address *model.ShippingAddress) *addressProto.ShippingAddress {
	return &addressProto.ShippingAddress{
		Id:              int64(address.Id),
		CustomerId:      int64(address.CustomerId),
		PostalCode:      address.PostalCode,
		CountryCode:     address.CountryCode,
		SubdivisionCode: address.SubdivisionCode,
		SubdivisionName: address.SubdivisionName,
		CityName:        address.CityName,
		AddressLine1:    address.AddressLine1,
		AddressLine2:    address.AddressLine2,
		Fullname:        address.Fullname,
		IsResidential:   address.IsResidential,
		RequestLiftgate: address.RequestLiftgate,
	}
}

func toBillingProto(address *model.BillingAddress) *addressProto.BillingAddress {
	return &addressProto.BillingAddress{
		Id:              int64(address.Id),
		CustomerId:      int64(address.CustomerId),
		PostalCode:      address.PostalCode,
		CountryCode:     address.CountryCode,
		SubdivisionCode: address.SubdivisionCode,
		SubdivisionName: address.SubdivisionName,
		CityName:        address.CityName,
		AddressLine1:    address.AddressLine1,
		AddressLine2:    address.AddressLine2,
		Fullname:        address.Fullname,
	}
}
