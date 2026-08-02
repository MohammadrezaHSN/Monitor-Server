package repository

import (
	"context"
	"time"

	"server/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RequestRepository struct {
	pool *pgxpool.Pool
}

func NewRequestRepository(pool *pgxpool.Pool) *RequestRepository {
	return &RequestRepository{pool: pool}
}

func (r *RequestRepository) Create(ctx context.Context, req *models.Request, agentInterval uint16) error {
	query := `
		INSERT INTO requests (
			agent_id,
			context,
			is_done,
			received_to_agent_at,
			estimate_response_at
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	if req.ReceivedToAgentAt != nil {
		estimate := req.ReceivedToAgentAt.Add(time.Duration(agentInterval) * time.Second)
		req.EstimateResponseAt = &estimate
	}

	err := r.pool.QueryRow(
		ctx,
		query,
		req.AgentID,
		req.Context,
		req.IsDone,
		req.ReceivedToAgentAt,
		req.EstimateResponseAt,
	).Scan(&req.ID, &req.CreatedAt)

	return err
}

func (r *RequestRepository) GetAllRequests(ctx context.Context) (*[]models.Request, error) {
	query := `
	SELECT * FROM requests
	`
	var requests []models.Request
	var request models.Request

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&request.ID, &request.CreatedAt)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return &requests, nil
}

func (r *RequestRepository) GetByID(ctx context.Context, id int64) (*models.Request, error) {
	query := `
		SELECT
			id,
			agent_id,
			context,
			is_done,
			created_at,
			received_to_agent_at,
			estimate_response_at
		FROM requests
		WHERE id = $1
	`

	var req models.Request

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&req.ID,
		&req.AgentID,
		&req.Context,
		&req.IsDone,
		&req.CreatedAt,
		&req.ReceivedToAgentAt,
		&req.EstimateResponseAt,
	)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func (r *RequestRepository) Update(ctx context.Context, req *models.Request, agentInterval uint16) error {
	query := `
		UPDATE requests
		SET
			agent_id = $1,
			context = $2,
			is_done = $3,
			received_to_agent_at = $4,
			estimate_response_at = $5
		WHERE id = $6
	`

	if req.ReceivedToAgentAt != nil {
		estimate := req.ReceivedToAgentAt.Add(time.Duration(agentInterval) * time.Second)
		req.EstimateResponseAt = &estimate
	} else {
		req.EstimateResponseAt = nil
	}

	_, err := r.pool.Exec(
		ctx,
		query,
		req.AgentID,
		req.Context,
		req.IsDone,
		req.ReceivedToAgentAt,
		req.EstimateResponseAt,
		req.ID,
	)

	return err
}

func (r *RequestRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM requests WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	return err
}
