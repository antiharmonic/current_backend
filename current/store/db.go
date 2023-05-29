package store


import (
	"context" // what is this
	"github.com/jackc/pgx/v4/pgxpool" //what is this
)

type postgres struct {
	pool *pgxpool.Pool
}

func CreatePostgresStore(connString string) (*postgres, error) {
	c, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	return &postgres{
		pool: c,
	}, nil
}