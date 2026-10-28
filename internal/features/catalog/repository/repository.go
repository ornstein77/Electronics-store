package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"kontakt/internal/core/domain"
	db "kontakt/internal/core/repository/postgres/pgx"
)

type Repository struct{ DB *pgxpool.Pool }

const columns = "id,name,spec,category,price,old,stock,type,color,badge"

func (r *Repository) List(ctx context.Context) ([]domain.Product, error) {
	rows, e := r.DB.Query(ctx, "SELECT "+columns+" FROM products ORDER BY name")
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []domain.Product{}
	for rows.Next() {
		var p domain.Product
		if e = rows.Scan(&p.ID, &p.Name, &p.Spec, &p.Category, &p.Price, &p.Old, &p.Stock, &p.Type, &p.Color, &p.Badge); e != nil {
			return nil, e
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
func (r *Repository) Get(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	var p domain.Product
	e := r.DB.QueryRow(ctx, "SELECT "+columns+" FROM products WHERE id=$1", id).Scan(&p.ID, &p.Name, &p.Spec, &p.Category, &p.Price, &p.Old, &p.Stock, &p.Type, &p.Color, &p.Badge)
	return p, db.Map(e)
}
func (r *Repository) Save(ctx context.Context, p domain.Product, create bool) error {
	if create {
		_, e := r.DB.Exec(ctx, "INSERT INTO products ("+columns+") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)", p.ID, p.Name, p.Spec, p.Category, p.Price, p.Old, p.Stock, p.Type, p.Color, p.Badge)
		return e
	}
	var id uuid.UUID
	return db.Map(r.DB.QueryRow(ctx, "UPDATE products SET name=$2,spec=$3,category=$4,price=$5,old=$6,stock=$7,type=$8,color=$9,badge=$10 WHERE id=$1 RETURNING id", p.ID, p.Name, p.Spec, p.Category, p.Price, p.Old, p.Stock, p.Type, p.Color, p.Badge).Scan(&id))
}
