package model

import (
	"database/sql"
	"time"
)

type Auth struct {
	ID         int64
	Info       AuthInfo
	created_at time.Time
	updated_at sql.NullTime
}

type AuthInfo struct {
	Name  string
	Email string
	Role  string
}
