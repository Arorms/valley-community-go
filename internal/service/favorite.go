package service

import (
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model/viewmodel"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/repository"
)

type FavoriteService interface {
	Toggle(userID string, postID int64) (bool, error)
	List(userID string, pageNo, pageSize int) ([]viewmodel.PostCard, int64, error)
}

type favoriteService struct {
	favorites repository.FavoriteRepository
	posts     repository.PostRepository
	cards     *cardBuilder
}

func NewFavoriteService(
	favorites repository.FavoriteRepository,
	posts repository.PostRepository,
	likes repository.LikeRepository,
	profiles repository.UserProfileRepository,
) FavoriteService {
	return &favoriteService{
		favorites: favorites,
		posts:     posts,
		cards:     newCardBuilder(likes, favorites, profiles),
	}
}

func (s *favoriteService) Toggle(userID string, postID int64) (bool, error) {
	if _, err := s.posts.GetByID(postID); err != nil {
		return false, err
	}
	added, err := s.favorites.Add(userID, postID)
	if err != nil {
		return false, err
	}
	if added {
		return true, nil
	}
	if _, err := s.favorites.Remove(userID, postID); err != nil {
		return false, err
	}
	return false, nil
}

func (s *favoriteService) List(userID string, pageNo, pageSize int) ([]viewmodel.PostCard, int64, error) {
	pageNo, pageSize = normalizePage(pageNo, pageSize)
	offset := (pageNo - 1) * pageSize

	posts, err := s.favorites.ListPosts(userID, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.favorites.Count(userID)
	if err != nil {
		return nil, 0, err
	}
	cards, err := s.cards.cards(posts, userID)
	if err != nil {
		return nil, 0, err
	}
	return cards, total, nil
}
