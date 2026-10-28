package pool

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func New(ctx context.Context, url string) (*pgxpool.Pool, error) {
	c, e := pgxpool.ParseConfig(url)
	if e != nil {
		return nil, e
	}
	c.MaxConns = 10
	c.MinConns = 1
	c.MaxConnLifetime = time.Hour
	p, e := pgxpool.NewWithConfig(ctx, c)
	if e != nil {
		return nil, e
	}
	if e = p.Ping(ctx); e != nil {
		p.Close()
		return nil, e
	}
	return p, nil
}
