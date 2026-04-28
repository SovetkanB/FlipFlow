package repo

import (
	"context"

	"github.com/SovetkanB/FlipFlow/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(dbCfg *config.DBConfig) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), dbCfg.DSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}

	return db, nil
}
