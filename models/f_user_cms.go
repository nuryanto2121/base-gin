package models

import uuid "github.com/satori/go.uuid"

type AddUserCms struct {
	Username        string      `json:"username" valid:"Required"`
	Password        string      `json:"password" valid:"Required"`
	ConfirmPassword string      `json:"confirm_password" valid:"Required"`
	Roles           []*AddRoles `json:"roles"`
}

type AddRoles struct {
	Role      string          `json:"role"`
	OutletIds []*AddOutletIds `json:"outlets"`
}
type AddOutletIds struct {
	OutletId uuid.UUID `json:"outlet_id"`
}
type ListUserCms struct {
	UserId   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
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
