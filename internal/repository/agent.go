package repository

import (
	"context"
	"time"

	"server/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AgentRepository struct {
	pool *pgxpool.Pool
}

func NewAgentRepository(pool *pgxpool.Pool) *AgentRepository {
	return &AgentRepository{pool: pool}
}

func (r *AgentRepository) Create(ctx context.Context, agent *models.Agent) error {
	query := `
		INSERT INTO agents (
			system_id,
			interval_seconds,
			is_alive,
			description,
			last_active,
			next_request_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	if agent.LastActive != nil {
		next := agent.LastActive.Add(time.Duration(agent.Interval) * time.Second)
		agent.NextRequestAt = &next
	}

	err := r.pool.QueryRow(
		ctx,
		query,
		agent.SystemID,
		agent.Interval,
		agent.IsAlive,
		agent.Description,
		agent.LastActive,
		agent.NextRequestAt,
	).Scan(&agent.ID, &agent.CreatedAt)

	return err
}

func (r *AgentRepository) GetAllAgents(ctx context.Context) (*[]models.Agent, error) {
	query := `
	SELECT * FROM agents
	`
	var agents []models.Agent
	var agent models.Agent

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&agent.ID, &agent.CreatedAt)
		if err != nil {
			return nil, err
		}
		agents = append(agents, agent)
	}
	return &agents, nil
}

func (r *AgentRepository) GetByID(ctx context.Context, id int64) (*models.Agent, error) {
	query := `
		SELECT
			id,
			system_id,
			interval_seconds,
			is_alive,
			description,
			last_active,
			created_at,
			next_request_at
		FROM agents
		WHERE id = $1
	`

	var agent models.Agent

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&agent.ID,
		&agent.SystemID,
		&agent.Interval,
		&agent.IsAlive,
		&agent.Description,
		&agent.LastActive,
		&agent.CreatedAt,
		&agent.NextRequestAt,
	)
	if err != nil {
		return nil, err
	}

	return &agent, nil
}

func (r *AgentRepository) Update(ctx context.Context, agent *models.Agent) error {
	query := `
		UPDATE agents
		SET
			system_id = $1,
			interval_seconds = $2,
			is_alive = $3,
			description = $4,
			last_active = $5,
			next_request_at = $6
		WHERE id = $7
	`

	if agent.LastActive != nil {
		next := agent.LastActive.Add(time.Duration(agent.Interval) * time.Second)
		agent.NextRequestAt = &next
	} else {
		agent.NextRequestAt = nil
	}

	_, err := r.pool.Exec(
		ctx,
		query,
		agent.SystemID,
		agent.Interval,
		agent.IsAlive,
		agent.Description,
		agent.LastActive,
		agent.NextRequestAt,
		agent.ID,
	)

	return err
}

func (r *AgentRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM agents WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	return err
}
