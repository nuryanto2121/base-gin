package models

type StatusOrder int64

const (
	SUBMITTED StatusOrder = iota
	APPROVE
	REJECT
)
