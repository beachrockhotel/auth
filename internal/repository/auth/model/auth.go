package model

import "database/sql"

import (
	"time"
)

type Role int32

type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	Role      string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
