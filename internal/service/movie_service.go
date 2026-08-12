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
	GetAllMovie(ctx context.Context) ([]*domain.Movie, error)
	GetByID(ctx context.Context, id int) (*domain.Movie, error)
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

func (s *MovieService) GetAllMovie(ctx context.Context) ([]*domain.Movie, error) {
	return s.repo.GetAllMovie(ctx)
}

func (s *MovieService) GetByID(ctx context.Context, id int) (*domain.Movie, error) {
	if id <= 0 {
		return &domain.Movie{}, fmt.Errorf("Id must be positive")
	}

	return s.repo.GetByID(ctx, id)
}