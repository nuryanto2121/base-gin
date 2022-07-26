package models

type OutletForm struct {
	OutletName   string             `json:"outlet_name" valid:"Required"`
	OutletCity   string             `json:"outlet_city" valid:"Required"`
	OutletDetail []*AddOutletDetail `json:"outlet_detail"`
}
