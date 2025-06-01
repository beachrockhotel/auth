package model

import (
	"database/sql"
	"time"
)

type Auth struct {
	ID        int64 `db:"id"`
	Info      *Info
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type Info struct {
	Name     string `db:"name"`
	Email    string `db:"email"`
	Role     string `db:"role"`
	Password string `db:"password"`
}
