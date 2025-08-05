package db

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitializeDatabase(pool *pgxpool.Pool, logger *slog.Logger) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			short_key VARCHAR(8) NOT NULL UNIQUE,
			original_url TEXT NOT NULL,
			user_id INTEGER REFERENCES users(id),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP + INTERVAL '24 hours'
		);
		CREATE TABLE IF NOT EXISTS download_history (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			url TEXT NOT NULL,
			format VARCHAR(4) NOT NULL,
			filename TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

	`)
	if err != nil {
		logger.Error("Failed to initialize database schema", "error", err)
		return err
	}
	logger.Info("Database schema initialized")
	return nil
}