package main

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4"
	"log"
	"time"
)

const (
	dbSN = "host=localhost port=54323 user=auth-user password=auth-password dbname=auth sslmode=disable"
)

func main() {
	ctx := context.Background()

	con, err := pgx.Connect(ctx, dbSN)
	if err != nil {
		log.Fatalf("failed to connect to database: #{err}")
	}

	defer con.Close(ctx)

	res, err := con.Exec(ctx, "INSERT INTO auth (name, email) VALUES ($1, $2)", gofakeit.Name(), gofakeit.Email())
	if err != nil {
		log.Fatalf("failed to insert data: #{err}")
	}

	log.Printf("inserted %d rows", res.RowsAffected())

	rows, err := con.Query(ctx, "SELECT id, name, email, role, created_at FROM auth")

	if err != nil {
		log.Fatalf("failed to select notes: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var name, email, role string
		var createdAt time.Time

		err = rows.Scan(&id, &name, &email, &role, &createdAt)
		if err != nil {
			log.Fatalf("failed to scan row: %v", err)
		}
		log.Printf("id: %d, name: %s, email: %s, role: %s, created_at: %v\n", id, name, email, role, createdAt)
	}
}
