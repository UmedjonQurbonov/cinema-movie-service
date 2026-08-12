package domain

import "time"

type Movie struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"`
	Age_limit   int       `json:"age_limit"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateMovieRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Age_limit   int    `json:"age_limit"`
}

type CreateMovieResponse struct {
	ID int `json:"id"`
}

type GetMovieResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Duration    int       `json:"duration"`
	Age_limit   int       `json:"age_limit"`
}
