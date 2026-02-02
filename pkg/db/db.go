package db

import (
	"context"
	"fmt"
	"log"
	"pi/internal/config"
	"pi/pkg/utils"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

func NewClient(ctx context.Context, maxAttempts int) (pool *pgxpool.Pool, err error) {
	var cfg config.ConfigDatabase
	err = cleanenv.ReadConfig("../../internal/config/db.yaml", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name))
		if err != nil {
			log.Fatal(err)
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	return
}
