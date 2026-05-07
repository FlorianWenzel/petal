package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type migration struct {
	name string
	sql  string
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, `create table if not exists schema_migrations (
		name text primary key,
		applied_at timestamptz not null default now()
	)`); err != nil {
		return err
	}

	migs, err := loadMigrations()
	if err != nil {
		return err
	}

	if err := backfillFromKysely(ctx, pool, migs); err != nil {
		return err
	}

	applied, err := loadApplied(ctx, pool)
	if err != nil {
		return err
	}

	for _, m := range migs {
		if applied[m.name] {
			continue
		}
		if err := applyMigration(ctx, pool, m); err != nil {
			return fmt.Errorf("migration %s: %w", m.name, err)
		}
	}
	return nil
}

func loadMigrations() ([]migration, error) {
	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return nil, err
	}
	var out []migration
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".up.sql") {
			continue
		}
		b, err := fs.ReadFile(migrationsFS, "migrations/"+e.Name())
		if err != nil {
			return nil, err
		}
		out = append(out, migration{
			name: strings.TrimSuffix(e.Name(), ".up.sql"),
			sql:  string(b),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return out, nil
}

func loadApplied(ctx context.Context, pool *pgxpool.Pool) (map[string]bool, error) {
	rows, err := pool.Query(ctx, `select name from schema_migrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := map[string]bool{}
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			return nil, err
		}
		out[n] = true
	}
	return out, rows.Err()
}

// If the previous Kysely-managed schema is present, copy its applied-migration
// names into schema_migrations so we don't re-run already-applied migrations.
// Anything Kysely never applied stays unmarked and runs as normal.
func backfillFromKysely(ctx context.Context, pool *pgxpool.Pool, _ []migration) error {
	var kyselyExists bool
	if err := pool.QueryRow(ctx,
		`select exists(select 1 from information_schema.tables where table_schema = 'public' and table_name = 'kysely_migration')`,
	).Scan(&kyselyExists); err != nil {
		return err
	}
	if !kyselyExists {
		return nil
	}

	rows, err := pool.Query(ctx, `select name from kysely_migration`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			return err
		}
		if _, err := pool.Exec(ctx,
			`insert into schema_migrations (name) values ($1) on conflict do nothing`,
			n,
		); err != nil {
			return err
		}
	}
	return rows.Err()
}

func applyMigration(ctx context.Context, pool *pgxpool.Pool, m migration) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, m.sql); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `insert into schema_migrations (name) values ($1)`, m.name); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
