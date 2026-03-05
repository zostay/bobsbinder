package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	migrate "github.com/rubenv/sql-migrate"
	"go.uber.org/zap"

	"github.com/zostay/bobsbinder/internal/config"
)

func Connect(cfg *config.Config, logger *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("connected to database",
		zap.String("host", cfg.DBHost),
		zap.String("database", cfg.DBName),
	)

	return db, nil
}

func RunMigrations(db *sql.DB, logger *zap.Logger) (int, error) {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	n, err := migrate.Exec(db, "mysql", migrations, migrate.Up)
	if err != nil {
		return 0, fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("migrations applied", zap.Int("count", n))
	return n, nil
}
