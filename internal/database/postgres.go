package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"booking-studio-foto/backend/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	SQL *sql.DB
}

func ConnectPostgres(cfg *config.Config) (*DB, error) {
	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(cfg.DatabaseMaxConns)
	db.SetMaxIdleConns(cfg.DatabaseMinConns)
	db.SetConnMaxIdleTime(cfg.DatabaseMaxIdle)
	db.SetConnMaxLifetime(cfg.DatabaseMaxLife)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Printf("[WARN] Database ping failed (will allow server to boot for healthcheck): %v", err)
		return &DB{SQL: db}, err
	}

	log.Println("[INFO] PostgreSQL connection pool established successfully")
	return &DB{SQL: db}, nil
}

func (d *DB) Close() error {
	if d.SQL != nil {
		return d.SQL.Close()
	}
	return nil
}
