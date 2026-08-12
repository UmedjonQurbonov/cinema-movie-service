package v1

import (
	"context"

	moviev1 "github.com/UmedjonQurbonov/cinema-libs/gen/go/movie/v1"
	"github.com/UmedjonQurbonov/cinema-movie-service/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"go.uber.org/zap"
)

type MovieService interface {
	Create(ctx context.Context, movie *domain.Movie) (int, error)
	GetAllMovie(ctx context.Context) ([]*domain.Movie, error)
	GetByID(ctx context.Context, id int) (*domain.Movie, error)
}

type MovieServer struct {
	moviev1.UnimplementedMovieServiceServer
	service MovieService
	log     *zap.Logger
}

func NewServer(service MovieService, log *zap.Logger) *MovieServer {
	return &MovieServer{
		service: service,
		log:     log,
	}
}

func (s *MovieServer) Create(ctx context.Context, req *moviev1.CreateMovieRequest) (*moviev1.CreateMovieResponse, error) {
	// 1. Формируем доменную модель с безопасными геттерами Protobuf
	movie := &domain.Movie{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Duration:    int(req.GetDuration()),
		Age_limit:   int(req.GetAgeLimit()),
	}

	// 2. Вызываем сервис
	id, err := s.service.Create(ctx, movie)
	if err != nil {
		s.log.Error("failed to create movie",
			zap.String("title", req.GetTitle()),
			zap.Error(err),
		)

		return nil, status.Error(codes.Internal, "internal error")
	}

	s.log.Info("movie created successfully",
		zap.Int64("id", int64(id)),
		zap.String("title", req.GetTitle()),
		zap.Int("duration_min", int(req.GetDuration())),
	)

	// 3. Возвращаем правильную структуру CreateMovieResponse (преобразовав int в int64)
	return &moviev1.CreateMovieResponse{
		Id: int64(id),
	}, nil
}		

func (s *MovieServer) List(ctx context.Context, req *moviev1.ListMovieRequest) (*moviev1.ListMovieResponse, error) {
	// 1. Получаем список доменных моделей из сервиса
	domainMovies, err := s.service.GetAllMovie(ctx)
	if err != nil {
		s.log.Error("failed to get movies", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}

	// 2. Аллоцируем срез под Protobuf-структуры с нужной емкостью (capacity)
	protoMovies := make([]*moviev1.MovieItem, 0, len(domainMovies))

	// 3. Перекладываем (маппим) данные из domain.Movie в moviev1.Movie
	for _, m := range domainMovies {
		protoMovies = append(protoMovies, &moviev1.MovieItem{
			Id:          int64(m.ID),
			Title:       m.Title,
			Duration:    int32(m.Duration),
			AgeLimit:    int32(m.Age_limit),
		})
	}

	s.log.Info("get movies successfully", zap.Int("count", len(protoMovies)))

	// 4. Оборачиваем срез в финальное сообщение ответа
	return &moviev1.ListMovieResponse{
		Movies: protoMovies,
	}, nil
}

func(s *MovieServer) GetById(ctx context.Context, req *moviev1.GetMovieRequest) (*moviev1.GetMovieResponse, error) {
	domainMovie, err := s.service.GetByID(ctx, int(req.Id))

	if err != nil {
		s.log.Error("failed to get movie", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	
	s.log.Info("get movies successfully", zap.Int("user:", int(req.Id)))

	return &moviev1.GetMovieResponse{
		Id: int64(domainMovie.ID),
		Title: domainMovie.Title,
		Description: domainMovie.Description,
		Duration: int32(domainMovie.Duration),
		AgeLimit: int32(domainMovie.Age_limit),
	}, nil
}