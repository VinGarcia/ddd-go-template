package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/vingarcia/ddd-go-template/v2-domain-adapters-and-helpers/adapters/log/jsonlogs"
	"github.com/vingarcia/ddd-go-template/v2-domain-adapters-and-helpers/domain"
	tt "github.com/vingarcia/ddd-go-template/v2-domain-adapters-and-helpers/helpers/testtools"
	migrations "github.com/vingarcia/ddd-go-template/v2-domain-adapters-and-helpers/migrations"
	"github.com/vingarcia/krest"
	"github.com/vingarcia/ksql"
	"github.com/vingarcia/ksql/adapters/kpgx"
	"golang.org/x/sync/errgroup"
)

type testData struct {
	db          ksql.Provider
	serverURL   string
	http        krest.Provider
	mockHandler *http.HandlerFunc
}

func TestAPI(t *testing.T) {
	ctx := context.Background()

	var port string
	getFreePorts(&port)

	dbURL, close := startPostgresDB(ctx, "ddd_go_template")
	defer close()

	db, err := kpgx.New(ctx, dbURL, ksql.Config{})
	tt.AssertNoErr(t, err)

	createDBSQL, err := migrations.Dir.ReadFile("create-db.sql")
	tt.AssertNoErr(t, err)
	_, err = db.Exec(ctx, string(createDBSQL))
	tt.AssertNoErr(t, err)

	restCli := krest.New(30 * time.Second)

	var mockHandler http.HandlerFunc

	// This fake server will replace any external server we depend on
	fakeServer := httptest.NewServer(&mockHandler)

	foursquareBaseURL := fakeServer.URL

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		err := startAPI(
			ctx,
			jsonlogs.New("INFO", domain.GetCtxValues),
			foursquareBaseURL,
			"fakeFoursquareClientID",
			"fakeFoursquareSecret",
			"", "", // Not using redis so we keep it with empty strings
			dbURL,
			port,
		)
		tt.AssertNoErr(t, err)
		return nil
	})

	time.Sleep(200 * time.Millisecond)

	testData := testData{
		db:          db,
		serverURL:   "http://localhost:" + port,
		http:        restCli,
		mockHandler: &mockHandler,
	}

	healthCheckTest(ctx, t, testData)
	usersTest(ctx, t, testData)

	cancel()
	g.Wait()
}

func healthCheckTest(ctx context.Context, t *testing.T, data testData) {
	resp, err := data.http.Get(ctx, data.serverURL, krest.RequestData{})
	tt.AssertNoErr(t, err)
	tt.AssertEqual(t, resp.StatusCode, 200)

	var dto struct {
		Service string `json:"service"`
		State   string `json:"state"`
	}
	err = json.Unmarshal(resp.Body, &dto)
	tt.AssertNoErr(t, err)

	tt.AssertEqual(t, dto.Service, "venues-service")
	tt.AssertEqual(t, dto.State, "healthy")
}

func getFreePorts(vars ...*string) error {
	for i := 0; i < len(vars); i++ {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		if err != nil {
			return err
		}

		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			return err
		}
		defer l.Close()
		*vars[i] = strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
	}

	return nil
}

func startPostgresDB(ctx context.Context, dbName string) (databaseURL string, closer func()) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not create a dockertest Pool: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "14.8",
			Env: []string{
				"POSTGRES_PASSWORD=postgres",
				"POSTGRES_USER=postgres",
				"POSTGRES_DB=" + dbName,
				"listen_addresses = '*'",
			},
		},
		func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		},
	)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://postgres:postgres@%s/%s?sslmode=disable", hostAndPort, dbName)

	fmt.Println("Connecting to postgres on url: ", databaseUrl)

	resource.Expire(40) // Tell docker to hard kill the container in 40 seconds

	var sqlDB *sql.DB
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 10 * time.Second
	pool.Retry(func() error {
		sqlDB, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return sqlDB.Ping()
	})
	if err != nil {
		log.Fatalf("Could not connect to postgres: %s", err)
	}

	sqlDB.Close()

	return databaseUrl, func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}
}

func resetTestState(ctx context.Context, data testData) {
	*data.mockHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte("mockHandler: unexpected request to: " + r.URL.Path))
	})

	data.db.Exec(ctx, "DELETE FROM users")
}
