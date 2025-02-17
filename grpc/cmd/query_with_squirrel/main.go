package main

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	dbDSN = "host=localhost port=54323 dbname=auth user=auth-user password=auth-password sslmode=disable"
)

func main() {
	ctx := context.Background()

	// Подключение к базе данных
	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Создание запроса на вставку
	builderInsert := squirrel.Insert("auth").
		PlaceholderFormat(squirrel.Dollar).
		Columns("name", "email").
		Values(gofakeit.Name(), gofakeit.Email()).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to generate insert query: %v", err)
	}

	var authID int
	err = pool.QueryRow(ctx, query, args...).Scan(&authID)
	if err != nil {
		log.Fatalf("failed to execute insert query: %v", err)
	}
	log.Printf("Inserted ID: %d", authID)

	// Запрос списка пользователей
	builderSelect := squirrel.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("auth").
		PlaceholderFormat(squirrel.Dollar).
		OrderBy("id ASC").
		Limit(10)

	query, args, err = builderSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to build select query: %v", err)
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to execute select query: %v", err)
	}
	defer rows.Close()

	var id int
	var name, email, role string
	var createdAt time.Time
	var updatedAt sql.NullTime

	for rows.Next() {
		err = rows.Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
		if err != nil {
			log.Fatalf("failed to scan row: %v", err)
		}
		log.Printf("id: #%d name: %s email: %s role: %s createdAt: %s, updatedAt: %s",
			id, name, email, role, createdAt, updatedAt.Time.String())
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("error iterating rows: %v", err)
	}

	// Обновление пользователя
	builderUpdate := squirrel.Update("auth").
		PlaceholderFormat(squirrel.Dollar).
		Set("name", gofakeit.Name()).
		Set("email", gofakeit.Email()).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": authID})

	query, args, err = builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to build update query: %v", err)
	}

	res, err := pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to execute update query: %v", err)
	}

	rowsAffected := res.RowsAffected()
	log.Printf("Updated rows: %d", rowsAffected)

	// Выборка конкретного пользователя
	builderSelectOne := squirrel.Select("id", "name", "email", "created_at", "updated_at").
		From("auth").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"id": authID}).
		Limit(1)

	query, args, err = builderSelectOne.ToSql()
	if err != nil {
		log.Fatalf("failed to build select query: %v", err)
	}

	err = pool.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &createdAt, &updatedAt)
	if err != nil {
		log.Fatalf("failed to execute select query: %v", err)
	}

	log.Printf("id: %d, name: %s, email: %s, created_at: %v, updated_at: %v",
		id, name, email, createdAt, updatedAt.Time)
}
