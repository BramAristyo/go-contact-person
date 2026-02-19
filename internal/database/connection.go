package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(dbUrl string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatal("Database connected unsuccessfully:", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("Unable to Ping Database!")
	}

	return pool
}
