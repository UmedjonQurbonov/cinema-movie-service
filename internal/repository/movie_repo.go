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

	if (err != nil) {
		return 0, fmt.Errorf("m.db.QueryRow: %w", err)
	}

	return id, nil
}