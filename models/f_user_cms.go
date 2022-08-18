package models

import uuid "github.com/satori/go.uuid"

type AddUserCms struct {
	Username        string          `json:"username" valid:"Required"`
	Name            string          `json:"name" valid:"Required"`
	PhoneNo         string          `json:"phone_no"`
	Email           string          `json:"email" `
	Password        string          `json:"password" valid:"Required"`
	ConfirmPassword string          `json:"confirm_password" valid:"Required"`
	Role            string          `json:"role" valid:"Required"`
	Outlets         []*AddOutletIds `json:"outlets"`
}

type EditUserCms struct {
	Username        string          `json:"username" valid:"Required"`
	Name            string          `json:"name" valid:"Required"`
	PhoneNo         string          `json:"phone_no"`
	Email           string          `json:"email" `
	Password        string          `json:"password" `
	ConfirmPassword string          `json:"confirm_password" `
	Role            string          `json:"role" valid:"Required"`
	Outlets         []*AddOutletIds `json:"outlets"`
}

type AddOutletIds struct {
	OutletId uuid.UUID `json:"outlet_id"`
}
type UserCms struct {
	UserId   uuid.UUID       `json:"user_id"`
	Username string          `json:"username"`
	Name     string          `json:"name"`
	Phone    string          `json:"phone_no"`
	Email    string          `json:"email"`
	Role     string          `json:"role"`
	RoleName string          `json:"role_name"`
	Outlest  []*OutletLookUp `json:"outlets" gorm:"-"`
}
type ResponseListUserCms struct {
	UserId   uuid.UUID       `json:"user_id"`
	Username string          `json:"username"`
	RoleCode []*UserRoleDesc `json:"group_code"`
}

type RoleCodes struct {
	RoleId   uuid.UUID       `json:"group_id"`
	RoleCode string          `json:"group_code"`
	Outlet   []*OutletLookUp `json:"outlets"`
}

type OutletLookUp struct {
	OutletId   uuid.UUID `json:"outlet_id"`
	OutletName string    `json:"outlet_name"`
}
