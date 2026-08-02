package repository

import (
	"context"

	"server/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FileRepository struct {
	pool *pgxpool.Pool
}

func NewFileRepository(pool *pgxpool.Pool) *FileRepository {
	return &FileRepository{pool: pool}
}

func (r *FileRepository) Create(ctx context.Context, file *models.File) error {
	query := `
		INSERT INTO files (
			agent_id,
			request_id,
			file_path_from_agent_system,
		    file_path_from_server_system,               
			is_completed,
			upload_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, received_to_server_at
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		file.AgentID,
		file.RequestID,
		file.FilePathFromAgentSystem,
		file.FilePathFromServerSystem,
		file.IsCompleted,
		file.UploadAt,
	).Scan(&file.ID, &file.ReceivedToServerAt)

	return err
}

func (r *FileRepository) GetByID(ctx context.Context, id int64) (*models.File, error) {
	query := `
		SELECT
			id,
			agent_id,
			request_id,
			file_path_from_agent_system,
		    file_path_from_server_system,        
			is_completed,
			upload_at,
			received_to_server_at
		FROM files
		WHERE id = $1
	`

	var file models.File

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&file.ID,
		&file.AgentID,
		&file.RequestID,
		&file.FilePathFromAgentSystem,
		&file.FilePathFromServerSystem,
		&file.IsCompleted,
		&file.UploadAt,
		&file.ReceivedToServerAt,
	)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (r *FileRepository) Update(ctx context.Context, file *models.File) error {
	query := `
		UPDATE files
		SET
			agent_id = $1,
			request_id = $2,
			file_path_from_agent_system = $3,
			file_path_from_server_system = $4,
			is_completed = $5,
			upload_at = $6
		WHERE id = $7
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		file.AgentID,
		file.RequestID,
		file.FilePathFromAgentSystem,
		file.FilePathFromServerSystem,
		file.IsCompleted,
		file.UploadAt,
		file.ID,
	)

	return err
}

func (r *FileRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM files WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *FileRepository) GetAllFiles(ctx context.Context) (*[]models.File, error) {
	query := `
	SELECT * FROM files
	`
	var files []models.File
	var file models.File

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&file.ID, &file.UploadAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return &files, nil
}
