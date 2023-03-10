package repo_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewStorageContiner() (func(context.Context) error, error) {
	ctx := context.Background()

	// Create the container
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "fsouza/fake-gcs-server",
			ExposedPorts: []string{
				"4443:4443",
			},
			Env: map[string]string{},
			Cmd: []string{
				"-scheme", "http",
			},
		},
	})
	if err != nil {
		return func(context.Context) error { return nil }, err
	}

	return c.Terminate, nil
}

// NewPostgresContainer creates a Postgres container and returns its DSN to be used
// in tests along with a termination callback to stop the container.
func NewPostgresContainer() (string, func(context.Context) error, error) {
	ctx := context.Background()

	templateURL := "postgres://postgres:postgres@localhost:%s/testdb?sslmode=disable&TimeZone=UTC"

	// Create the container
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "postgres:14.1",
			ExposedPorts: []string{
				"0:5432",
			},
			Env: map[string]string{
				"POSTGRES_DB":       "testdb",
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
				"POSTGRES_SSL_MODE": "disable",
			},
			Cmd: []string{
				"postgres", "-c", "fsync=off",
			},
			WaitingFor: wait.ForSQL(
				"5432/tcp",
				"postgres",
				func(p nat.Port) string {
					return fmt.Sprintf(templateURL, p.Port())
				},
			).Timeout(time.Second * 30),
		},
	})
	if err != nil {
		return "", func(context.Context) error { return nil }, err
	}

	// Find ports assigned to the new container
	ports, err := c.Ports(ctx)
	if err != nil {
		return "", func(context.Context) error { return nil }, err
	}

	// Format driverURL
	driverURL := fmt.Sprintf(templateURL, ports["5432/tcp"][0].HostPort)

	return driverURL, c.Terminate, nil
}

func migrateUp(db *sql.DB, files string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	// find out the absolute path to this file
	// it'll be used to determine the project's root path
	_, callerPath, _, _ := runtime.Caller(0)

	// look for migrations source starting from project's root dir
	sourceURL := fmt.Sprintf(
		"file://%s/../../../%s",
		filepath.ToSlash(filepath.Dir(callerPath)),
		filepath.ToSlash(files),
	)

	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
