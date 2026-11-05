package service

import (
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/content"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model/viewmodel"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/repository"
)

type cardBuilder struct {
	likes     repository.LikeRepository
	favorites repository.FavoriteRepository
	profiles  repository.UserProfileRepository
}

func newCardBuilder(likes repository.LikeRepository, favorites repository.FavoriteRepository, profiles repository.UserProfileRepository) *cardBuilder {
	return &cardBuilder{likes: likes, favorites: favorites, profiles: profiles}
}

func (b *cardBuilder) cards(posts []model.Post, viewerID string) ([]viewmodel.PostCard, error) {
	if len(posts) == 0 {
		return []viewmodel.PostCard{}, nil
	}

	authorIDs := make([]string, 0, len(posts))
	postIDs := make([]int64, 0, len(posts))
	for _, post := range posts {
		authorIDs = append(authorIDs, post.AuthorID)
		postIDs = append(postIDs, post.ID)
	}

	profiles, err := b.profiles.ListByIDs(authorIDs)
	if err != nil {
		return nil, err
	}

	liked := map[int64]bool{}
	favorited := map[int64]bool{}
	if viewerID != "" {
		if liked, err = b.likes.LikedPosts(viewerID, postIDs); err != nil {
			return nil, err
		}
		if favorited, err = b.favorites.FavoritedPosts(viewerID, postIDs); err != nil {
			return nil, err
		}
	}

	cards := make([]viewmodel.PostCard, 0, len(posts))
	for _, post := range posts {
		cards = append(cards, viewmodel.PostCard{
			ID:           post.ID,
			Title:        post.Title,
			Excerpt:      content.Excerpt(post.MarkdownContent, 160),
			Author:       viewmodel.NewUserSummary(post.AuthorID, profilePointer(profiles, post.AuthorID)),
			CreatedAt:    post.CreatedAt,
			LikeCount:    post.LikeCount,
			CommentCount: post.CommentCount,
			Liked:        liked[post.ID],
			Favorited:    favorited[post.ID],
		})
	}
	return cards, nil
}

func profilePointer(profiles map[string]model.UserProfile, userID string) *model.UserProfile {
	profile, ok := profiles[userID]
	if !ok {
		return nil
	}
	return &profile
}
