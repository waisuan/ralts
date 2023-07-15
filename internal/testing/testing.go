package testing

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"ralts/internal/config"
)

var ctx = context.Background()

func InitTestResources(cfg *config.Config) func() {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisConn,
		DB:   0, // use default DB
	})

	rdb.FlushAll(ctx)
	ClearDB(cfg)
	InitDB(cfg)

	return func() {
		rdb.FlushAll(ctx)
		ClearDB(cfg)
	}
}

func InitDB(cfg *config.Config) {
	m, err := migrate.New(
		"file:../../db/migrations",
		cfg.DatabaseConn)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Warn(err)
	}
}

func ClearDB(cfg *config.Config) {
	m, err := migrate.New(
		"file:../../db/migrations",
		cfg.DatabaseConn)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Down(); err != nil {
		log.Warn(err)
	}
}

//func TestContainers() {
//	ctx := context.Background()
//
//	containerReq := testcontainers.ContainerRequest{
//		Image:        "postgres:latest",
//		ExposedPorts: []string{"5432/tcp"},
//		WaitingFor:   wait.ForListeningPort("5432/tcp"),
//		Env: map[string]string{
//			"POSTGRES_DB":       "ralts_test",
//			"POSTGRES_PASSWORD": "postgres",
//			"POSTGRES_USER":     "postgres",
//		},
//	}
//
//	dbContainer, _ := testcontainers.GenericContainer(
//		context.Background(),
//		testcontainers.GenericContainerRequest{
//			ContainerRequest: containerReq,
//			Started:          true,
//		})
//
//	host, _ := dbContainer.Host(context.Background())
//	port, _ := dbContainer.MappedPort(context.Background(), "5432")
//
//	redisContainer, err := redis.RunContainer(ctx)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	rawRedisConn, _ := redisContainer.ConnectionString(ctx)
//	redisConn, found := strings.CutPrefix(rawRedisConn, "redis://")
//
//	if err := dbContainer.Terminate(ctx); err != nil {
//		t.Fatalf("failed to terminate pgContainer: %s", err)
//	}
//
//	if err := redisContainer.Terminate(ctx); err != nil {
//		t.Fatalf("failed to terminate container: %s", err)
//	}
//}
