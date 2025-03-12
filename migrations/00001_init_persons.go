package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upPersons, downPersons)
}

func upPersons(ctx context.Context, tx *sql.Tx) error {
	query := `CREATE TABLE IF NOT EXISTS persons (
    	id UUID PRIMARY KEY NOT NULL,
    	email TEXT NOT NULL UNIQUE,
    	login VARCHAR(40) NOT NULL,
    	about_me TEXT NOT NULL,
    	password_hash TEXT NOT NULL,
    	image BYTEA
	)`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}

func downPersons(ctx context.Context, tx *sql.Tx) error {
	query := `DROP TABLE IF EXISTS persons`
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}
