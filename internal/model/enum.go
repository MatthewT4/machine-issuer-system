package model

type UserStatus = int

const (
	ACTIVE UserStatus = iota
	DELETED
)
