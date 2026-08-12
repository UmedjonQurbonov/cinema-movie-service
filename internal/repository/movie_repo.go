package repository

import (
	"context"
	"fmt"

	"github.com/UmedjonQurbonov/cinema-movie-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MovieRepo struct {
	db *pgxpool.Pool
}

func NewMovieRepo(db *pgxpool.Pool) *MovieRepo {
	return &MovieRepo{
		db: db,
	}
}

func (r *MovieRepo) Create(ctx context.Context, movie *domain.Movie) (int, error) {

	const query = "INSERT INTO movies (title, description, duration, age_limit) VALUES ($1, $2, $3, $4) RETURNING id"

	var id int

	err := r.db.QueryRow(ctx, query, movie.Title, movie.Description, movie.Duration, movie.Age_limit).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("m.db.QueryRow: %w", err)
	}

	return id, nil
}

func (r *MovieRepo) GetAllMovie(ctx context.Context) ([]*domain.Movie, error) {
	query := `SELECT id, title, description, duration, age_limit FROM movies ORDER BY id DESC`
	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("repo.GetAll execute query: %w", err)
	}

	defer rows.Close()

	movies := make([]*domain.Movie, 0)

	for rows.Next() {
		var m domain.Movie

		// 4. Сканируем значения колонок строго в том порядке, в котором указали в SELECT
		err := rows.Scan(
			&m.ID,
			&m.Title,
			&m.Description,
			&m.Duration,
			&m.Age_limit,
		)
		if err != nil {
			return nil, fmt.Errorf("repo.GetAll scan row: %w", err)
		}

		movies = append(movies, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo.GetAll rows iteration: %w", err)
	}

	return movies, nil
}


func (r *MovieRepo) GetByID(ctx context.Context, id int) (*domain.Movie, error) {
	const query = "SELECT id, title, description, duration, age_limit, created_at FROM movies WHERE id = $1" 

	var m domain.Movie
	err := r.db.QueryRow(ctx, query, id).Scan(
		&m.ID,
		&m.Title,
		&m.Description,
		&m.Duration,
		&m.Age_limit,
		&m.CreatedAt,
	)
	if err != nil {
		return &domain.Movie{}, fmt.Errorf("r.db.QueryRow: %w", err)
	}

	return &m, nil
}