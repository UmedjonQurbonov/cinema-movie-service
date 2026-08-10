package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/UmedjonQurbonov/cinema-movie-service/internal/domain"
)

var (
	ErrMovieNotFound = errors.New("movie not found")
	ErrInvalidData   = errors.New("invalid movie data")
)

type MovieRepository interface {
	Create(ctx context.Context, movie *domain.Movie) (int, error)
}


type MovieService struct {
	repo MovieRepository
}

func NewMovieService(repo MovieRepository) *MovieService {
	return &MovieService{
		repo: repo,
	}
}

func (s *MovieService) Create(ctx context.Context, movie *domain.Movie) (int, error) {
	
	if strings.TrimSpace(movie.Title) == "" {
		return 0, fmt.Errorf("%w: title cannot be empty", ErrInvalidData)
	} 

	if strings.TrimSpace(movie.Title) == "" {
		return 0, fmt.Errorf("%w: description cannot be empty", ErrInvalidData)
	} 

	if movie.Duration < 0 {
		return 0, fmt.Errorf("%w: Duration must be positive", ErrInvalidData)
	}

	if movie.Duration < 0 {
		return 0, fmt.Errorf("%w: Ale limit must be positive", ErrInvalidData)
	}

	return s.repo.Create(ctx, movie)
}