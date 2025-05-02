package model

import (
	"database/sql"
	"time"
)

type Auth struct {
	ID        int64
	Info      AuthInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type AuthInfo struct {
	Title   string
	Content string
}
