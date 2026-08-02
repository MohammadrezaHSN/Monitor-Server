package db

import (
	"context"
	"log"
	"server/internal/repository"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

/*

For Migrate, Must Run:
docker run --rm -v ${PWD}/migrations:/migrations migrate/migrate -path=/migrations/ -database "postgres://user:password@host.docker.internal:5432/mydatabase?sslmode=disable" up

For Remove Migrations Run (2 at the end of below command it means remove 2 last migrations.):
docker run --rm -v ${PWD}/migrations:/migrations migrate/migrate -path=/migrations/ -database "postgres://user:password@host.docker.internal:5432/mydatabase?sslmode=disable" down 2

For Open PostgresSql On Terminal Must Run:
docker exec -it my-go-postgres psql -U user -d mydatabase

For Shows List Of Tables:
	\dt

For Shows Each Table Run:
	\d table_name => \d agents

	OR

	SELECT * FROM agents;

For Exit Run:
	\q

*/

func NewPostgresPool(connStr string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnIdleTime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

type RepoQueries struct {
	AgentRepo    *repository.AgentRepository
	RequestRepo  *repository.RequestRepository
	ResponseRepo *repository.ResponseRepository
	FileRepo     *repository.FileRepository
}

func PrepareQueries() (RepoQueries, context.Context) {
	connStr := "postgres://user:password@localhost:5432/mydatabase?sslmode=disable"
	pool, err := NewPostgresPool(connStr)
	if err != nil {
		log.Printf("Failed to initialize database pool: %v\n", err)
		connStr = "postgres://user:password@host.docker.internal:5432/mydatabase?sslmode=disable"
		err = nil
		pool, err = NewPostgresPool(connStr)
		if err != nil {
			log.Fatalf("Failed to initialize database pool: %v\n", err)
		}
	}
	defer func() {
		log.Println("Closing database pool...")
		pool.Close()
	}()

	log.Println("Connection Path is: ", connStr)

	log.Println("Database connection pool established successfully")
	var rq RepoQueries

	// 3. Initialize repositories
	rq.AgentRepo = repository.NewAgentRepository(pool)
	rq.RequestRepo = repository.NewRequestRepository(pool)
	rq.ResponseRepo = repository.NewResponseRepository(pool)
	rq.FileRepo = repository.NewFileRepository(pool)

	ctx := context.Background()
	return rq, ctx
}
