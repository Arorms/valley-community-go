package services

import (
	"valley-community-go/models"
	"valley-community-go/repositories"
)

type PostService interface {
	Create(request models.PostCreateRequest) (*models.Post, error)
	GetById(id uint) (*models.Post, error)
	List(offset int, limit int) ([]models.Post, error)
}

type postService struct {
	repository repositories.PostRepository
}

func NewPostService(repository repositories.PostRepository) PostService {
	return &postService{repository: repository}
}

func (s *postService) Create(request models.PostCreateRequest) (*models.Post, error) {
	post := &models.Post{
		Title:   request.Title,
		Content: request.Content,
	}
	if error := s.repository.Create(post); error != nil {
		return nil, error
	}
	return post, nil
}

func (s *postService) GetById(id uint) (*models.Post, error) {
	return s.repository.GetById(id)
}

func (s *postService) List(offset int, limit int) ([]models.Post, error) {
	return s.repository.List(offset, limit)
}
