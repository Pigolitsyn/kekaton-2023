package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"kekaton/back/internal/config"
	"kekaton/back/internal/server"
	"kekaton/back/internal/service"
	"kekaton/back/internal/storage"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		panic(err)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errc := make(chan error, 1)

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	db, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		return err
	}

	stg := storage.New(db, storage.Config{
		Timeout: time.Second,
	})

	src := service.New(stg)

	app := fiber.New(fiber.Config{
		AppName:      cfg.Name,
		ErrorHandler: server.ErrorHandler,
	})

	srv := server.New(ctx, app, src, server.Config{
		Secret:    []byte(cfg.Secret),
		TokenName: "jwt",
	})

	go srv.Listen(cfg.Port, errc)

	select {
	case err = <-errc:
		return err
	case <-ctx.Done():
		return srv.Shutdown()
	}
}
