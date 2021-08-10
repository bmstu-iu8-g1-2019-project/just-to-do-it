package infrastructure

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"just-to-do-it/internal/config"
	"just-to-do-it/internal/interfaces"
)

type PostgresClient struct {
	Pool *pgxpool.Pool
	Ctx  context.Context
}

func initPostgresClient(cfg *config.Config, ctx context.Context) (interfaces.DBHandler, error) {
	dbPool, err := pgxpool.Connect(ctx, fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresUsername, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDatabase))
	if err != nil {
		return nil, errors.Wrap(err, "postgres init")
	}
	return &PostgresClient{Pool: dbPool, Ctx: ctx}, nil
}

func (p *PostgresClient) GetPool() *pgxpool.Pool {
	return p.Pool
}

func (p *PostgresClient) AcquireConn() (*pgxpool.Conn, error) {
	return p.Pool.Acquire(p.Ctx)
}

func (p *PostgresClient) GetCtx() context.Context {
	return p.Ctx
}
