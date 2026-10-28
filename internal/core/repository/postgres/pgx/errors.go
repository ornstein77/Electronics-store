package pgx

import (
	"errors"
	"github.com/jackc/pgx/v5"
	core "kontakt/internal/core/errors"
)

func Map(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return core.NotFound
	}
	return err
}
