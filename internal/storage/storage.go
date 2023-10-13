package storage

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Config struct {
	Timeout time.Duration
}

type Storage struct {
	pool   *pgxpool.Pool
	config Config
}

func New(pool *pgxpool.Pool, config Config) *Storage {
	return &Storage{
		pool:   pool,
		config: config,
	}
}
