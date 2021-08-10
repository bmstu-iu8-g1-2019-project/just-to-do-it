package interfaces

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DBHandler interface {
	GetPool() *pgxpool.Pool
	AcquireConn() (*pgxpool.Conn, error)
	GetCtx() context.Context
}
