package models

import uuid "github.com/satori/go.uuid"

type AddUserCms struct {
	Username        string         `json:"username" valid:"Required"`
	Password        string         `json:"password" valid:"Required"`
	ConfirmPassword string         `json:"confirm_password" valid:"Required"`
	Groups          []*AddGroupIds `json:"groups"`
}

type AddGroupIds struct {
	GroupId   uuid.UUID       `json:"group_id"`
	OutletIds []*AddOutletIds `json:"outlets"`
}
type AddOutletIds struct {
	OutletId uuid.UUID `json:"outlet_id"`
}

type ListUserCms struct {
	UserId      uuid.UUID `json:"user_id"`
	Username    string    `json:"username"`
	GroupCode   string    `json:"group_code"`
	GroupOutlet string    `json:"group_outlet" `
}
