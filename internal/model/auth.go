package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID         int64
	Info       AuthInfo
	Created_at time.Time
	Updated_at sql.NullTime
}

type AuthInfo struct {
	Name  string
	Email string
	Role  string
}
