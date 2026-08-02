package repository

import (
	"context"
	"server/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ResponseRepository struct {
	pool *pgxpool.Pool
}

func NewResponseRepository(pool *pgxpool.Pool) *ResponseRepository {
	return &ResponseRepository{pool: pool}
}

func (r *ResponseRepository) Create(ctx context.Context, resp *models.Response) error {
	query := `
		INSERT INTO responses (
			agent_id,
			request_id,
			result,
			is_successful,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, received_to_server_at
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		resp.AgentID,
		resp.RequestID,
		resp.Result,
		resp.IsSuccessful,
		resp.CreatedAt,
	).Scan(&resp.ID, &resp.ReceivedToServerAt)

	return err
}

func (r *ResponseRepository) GetAllResponses(ctx context.Context) (*[]models.Response, error) {
	query := `
	SELECT * FROM response
	`
	var responses []models.Response
	var response models.Response

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&response.ID, &response.CreatedAt)
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}
	return &responses, nil

}

func (r *ResponseRepository) GetByID(ctx context.Context, id int64) (*models.Response, error) {
	query := `
		SELECT
			id,
			agent_id,
			request_id,
			result,
			is_successful,
			created_at,
			received_to_server_at
		FROM responses
		WHERE id = $1
	`

	var resp models.Response

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&resp.ID,
		&resp.AgentID,
		&resp.RequestID,
		&resp.Result,
		&resp.IsSuccessful,
		&resp.CreatedAt,
		&resp.ReceivedToServerAt,
	)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (r *ResponseRepository) Update(ctx context.Context, resp *models.Response) error {
	query := `
		UPDATE responses
		SET
			agent_id = $1,
			request_id = $2,
			result = $3,
			is_successful = $4,
			created_at = $5
		WHERE id = $6
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		resp.AgentID,
		resp.RequestID,
		resp.Result,
		resp.IsSuccessful,
		resp.CreatedAt,
		resp.ID,
	)

	return err
}

func (r *ResponseRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM responses WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	return err
}
