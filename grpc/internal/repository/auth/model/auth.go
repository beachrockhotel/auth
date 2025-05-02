package model

import (
	"time"
)

type Role int32

type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	Role      Role
	CreatedAt time.Time
	UpdatedAt *time.Time
}
