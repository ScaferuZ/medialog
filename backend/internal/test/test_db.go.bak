package testutil

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/medialogg/backend/internal/db"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDatabase struct {
	Container   *postgres.PostgresContainer
	Conn        *pgx.Conn
	Queries     *db.Queries
	DatabaseURL string
}

func NewTestDatabase(ctx context.Context) (*TestDatabase, error) {
	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("medialogg_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2)),
	)
	if err != nil {
		return nil, fmt.Errorf("start postgres test container: %w", err)
	}

	databaseURL, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("get test database connection string: %w", err)
	}

	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("connect to test database: %w", err)
	}

	testDB := &TestDatabase{
		Container:   container,
		Conn:        conn,
		Queries:     db.New(conn),
		DatabaseURL: databaseURL,
	}

	if err := testDB.applyMigrations(ctx); err != nil {
		_ = testDB.Close(ctx)
		return nil, err
	}

	return testDB, nil
}

func (tdb *TestDatabase) Reset(ctx context.Context) error {
	_, err := tdb.Conn.Exec(ctx, "TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	if err != nil {
		return fmt.Errorf("reset test database: %w", err)
	}

	return nil
}

func (tdb *TestDatabase) Close(ctx context.Context) error {
	var errs []error

	if tdb.Conn != nil {
		errs = append(errs, tdb.Conn.Close(ctx))
	}

	if tdb.Container != nil {
		errs = append(errs, tdb.Container.Terminate(ctx))
	}

	return errors.Join(errs...)
}

func (tdb *TestDatabase) applyMigrations(ctx context.Context) error {
	migrationsDir, err := migrationsPath()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations directory: %w", err)
	}

	filenames := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}

		filenames = append(filenames, entry.Name())
	}

	sort.Strings(filenames)

	for _, filename := range filenames {
		migrationPath := filepath.Join(migrationsDir, filename)
		migrationSQL, err := os.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", filename, err)
		}

		statement := strings.TrimSpace(string(migrationSQL))
		if statement == "" {
			continue
		}

		if _, err := tdb.Conn.Exec(ctx, statement); err != nil {
			return fmt.Errorf("apply migration %s: %w", filename, err)
		}
	}

	return nil
}

func migrationsPath() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("resolve test database file path")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", "..", "migrations")), nil
}
