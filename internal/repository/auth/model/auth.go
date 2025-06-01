package model

import (
	"database/sql"
	"github.com/beachrockhotel/auth/pkg/auth_v1"
	"time"
)

type Auth struct {
	ID        int64
	Info      AuthInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type AuthInfo struct {
	Name  string
	Email string
	Role  auth_v1.Role
}
