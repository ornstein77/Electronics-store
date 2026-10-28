package repository

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"kontakt/internal/core/domain"
	core "kontakt/internal/core/errors"
	db "kontakt/internal/core/repository/postgres/pgx"
	"sort"
)

type Repository struct{ DB *pgxpool.Pool }

func (r *Repository) Create(ctx context.Context, in domain.OrderInput) (domain.Order, error) {
	o := domain.Order{ID: uuid.New(), OrderInput: in, Status: "Новый"}
	tx, e := r.DB.Begin(ctx)
	if e != nil {
		return o, e
	}
	defer tx.Rollback(ctx)
	// Stable row lock order prevents deadlocks between simultaneous carts.
	sort.Slice(o.Items, func(i, j int) bool { return o.Items[i].ID.String() < o.Items[j].ID.String() })
	for i := range o.Items {
		it := &o.Items[i]
		var stock int
		e = tx.QueryRow(ctx, "SELECT name,price,stock FROM products WHERE id=$1 FOR UPDATE", it.ID).Scan(&it.Name, &it.Price, &stock)
		if e != nil {
			return o, db.Map(e)
		}
		if stock < it.Qty {
			return o, core.Conflict
		}
		o.Total += it.Price * int64(it.Qty)
		if _, e = tx.Exec(ctx, "UPDATE products SET stock=stock-$2 WHERE id=$1", it.ID, it.Qty); e != nil {
			return o, e
		}
	}
	body, e := json.Marshal(o.OrderInput)
	if e != nil {
		return o, e
	}
	e = tx.QueryRow(ctx, "INSERT INTO orders(id,body,total,status) VALUES($1,$2,$3,$4) RETURNING created_at", o.ID, body, o.Total, o.Status).Scan(&o.CreatedAt)
	if e != nil {
		return o, e
	}
	return o, tx.Commit(ctx)
}
func (r *Repository) List(ctx context.Context) ([]domain.Order, error) {
	rows, e := r.DB.Query(ctx, "SELECT id,body,total,status,created_at FROM orders ORDER BY created_at DESC LIMIT 500")
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []domain.Order{}
	for rows.Next() {
		var o domain.Order
		var b []byte
		if e = rows.Scan(&o.ID, &b, &o.Total, &o.Status, &o.CreatedAt); e != nil {
			return nil, e
		}
		if e = json.Unmarshal(b, &o.OrderInput); e != nil {
			return nil, e
		}
		out = append(out, o)
	}
	return out, rows.Err()
}
func (r *Repository) Status(ctx context.Context, id uuid.UUID, status string) error {
	var result uuid.UUID
	return db.Map(r.DB.QueryRow(ctx, "UPDATE orders SET status=$2 WHERE id=$1 RETURNING id", id, status).Scan(&result))
}
