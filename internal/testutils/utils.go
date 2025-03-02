package testutils

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
)

var user = "jed"
var password = "jed"
var databaseName = "jed_shop"

type DBResources struct {
	Container *gnomock.Container
	Pool      *pgxpool.Pool
}

// NewDBResources creates a new DBResources instance with a new gnomock container and pool.
func NewDBResources(t *testing.T, queries []string) *DBResources {
	t.Helper()
	container := newDBContainer(t, queries)
	pool := newDBPool(t, container)
	return &DBResources{Container: container, Pool: pool}
}

func (dbr *DBResources) Close(t *testing.T) {
	t.Helper()
	dbr.Pool.Close()
	gnomock.Stop(dbr.Container)
}

func newDBContainer(t *testing.T, queries []string) *gnomock.Container {
	t.Helper()
	queries = append([]string{"CREATE EXTENSION IF NOT EXISTS pgcrypto;"}, queries...)
	p := postgres.Preset(
		postgres.WithUser(user, password),
		postgres.WithDatabase(databaseName),
		postgres.WithQueries(queries...),
		postgres.WithQueriesFile("../../../db/schema.sql"),
		postgres.WithTimezone("Europe/Belgrade"),
	)
	container, err := gnomock.Start(p)
	if err != nil {
		t.Fatal(err)
	}
	return container
}

func newDBPool(t *testing.T, container *gnomock.Container) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()
	port := container.DefaultPort()
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, container.Host, port, databaseName)
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		t.Fatal(err)
	}
	return pool
}
