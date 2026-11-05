package service

import (
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/repository"
)

type PostService interface {
	Create(request model.PostCreateRequest) (*model.Post, error)
	GetById(id int64) (*model.Post, error)
	List(offset int, limit int) ([]model.Post, error)
}

type postService struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) PostService {
	return &postService{repo: repo}
}

func (s *postService) Create(request model.PostCreateRequest) (*model.Post, error) {
	post := &model.Post{
		Title:   request.Title,
		Content: request.Content,
	}
	if error := s.repo.Create(post); error != nil {
		return nil, error
	}
	return post, nil
}

func (s *postService) GetById(id int64) (*model.Post, error) {
	return s.repo.GetById(id)
}

func (s *postService) List(offset int, limit int) ([]model.Post, error) {
	return s.repo.List(offset, limit)
}
