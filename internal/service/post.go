package service

import (
	"errors"
	"strings"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/content"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model/request"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model/viewmodel"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/repository"
)

const PageSize = 20

type PostService interface {
	Create(authorID string, req request.PostCreateRequest) (*model.Post, error)
	GetByID(id int64) (*model.Post, error)
	List(pageNo, pageSize int, viewerID string) ([]viewmodel.PostCard, int64, error)
	ListByAuthor(authorID string, pageNo, pageSize int, viewerID string) ([]viewmodel.PostCard, int64, error)
	Detail(postID int64, viewerID string) (*viewmodel.PostDetail, []viewmodel.CommentView, error)
}

type postService struct {
	posts     repository.PostRepository
	comments  repository.CommentRepository
	likes     repository.LikeRepository
	favorites repository.FavoriteRepository
	profiles  repository.UserProfileRepository
	cards     *cardBuilder
}

func NewPostService(
	posts repository.PostRepository,
	comments repository.CommentRepository,
	likes repository.LikeRepository,
	favorites repository.FavoriteRepository,
	profiles repository.UserProfileRepository,
) PostService {
	return &postService{
		posts:     posts,
		comments:  comments,
		likes:     likes,
		favorites: favorites,
		profiles:  profiles,
		cards:     newCardBuilder(likes, favorites, profiles),
	}
}

func (s *postService) Create(authorID string, req request.PostCreateRequest) (*model.Post, error) {
	title := strings.TrimSpace(req.Title)
	markdown := strings.TrimSpace(req.MarkdownContent)
	if title == "" {
		return nil, errors.New("title is required")
	}
	if len([]rune(title)) > 200 {
		return nil, errors.New("title is too long")
	}
	if markdown == "" {
		return nil, errors.New("content is required")
	}

	html, err := content.Render(markdown)
	if err != nil {
		return nil, err
	}

	post := &model.Post{
		AuthorID:        authorID,
		Title:           title,
		MarkdownContent: markdown,
		HTMLContent:     html,
	}
	if err := s.posts.Create(post); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *postService) GetByID(id int64) (*model.Post, error) {
	return s.posts.GetByID(id)
}

func (s *postService) List(pageNo, pageSize int, viewerID string) ([]viewmodel.PostCard, int64, error) {
	pageNo, pageSize = normalizePage(pageNo, pageSize)
	offset := (pageNo - 1) * pageSize

	posts, err := s.posts.List(offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.posts.Count()
	if err != nil {
		return nil, 0, err
	}
	cards, err := s.cards.cards(posts, viewerID)
	if err != nil {
		return nil, 0, err
	}
	return cards, total, nil
}

func (s *postService) ListByAuthor(authorID string, pageNo, pageSize int, viewerID string) ([]viewmodel.PostCard, int64, error) {
	pageNo, pageSize = normalizePage(pageNo, pageSize)
	offset := (pageNo - 1) * pageSize

	posts, err := s.posts.ListByAuthor(authorID, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.posts.CountByAuthor(authorID)
	if err != nil {
		return nil, 0, err
	}
	cards, err := s.cards.cards(posts, viewerID)
	if err != nil {
		return nil, 0, err
	}
	return cards, total, nil
}

func (s *postService) Detail(postID int64, viewerID string) (*viewmodel.PostDetail, []viewmodel.CommentView, error) {
	post, err := s.posts.GetByID(postID)
	if err != nil {
		return nil, nil, err
	}

	profiles, err := s.profiles.ListByIDs([]string{post.AuthorID})
	if err != nil {
		return nil, nil, err
	}

	liked := false
	favorited := false
	if viewerID != "" {
		likedPosts, err := s.likes.LikedPosts(viewerID, []int64{postID})
		if err != nil {
			return nil, nil, err
		}
		liked = likedPosts[postID]

		favoritedPosts, err := s.favorites.FavoritedPosts(viewerID, []int64{postID})
		if err != nil {
			return nil, nil, err
		}
		favorited = favoritedPosts[postID]
	}

	detail := &viewmodel.PostDetail{
		PostCard: viewmodel.PostCard{
			ID:           post.ID,
			Title:        post.Title,
			Author:       viewmodel.NewUserSummary(post.AuthorID, profilePointer(profiles, post.AuthorID)),
			CreatedAt:    post.CreatedAt,
			LikeCount:    post.LikeCount,
			CommentCount: post.CommentCount,
			Liked:        liked,
			Favorited:    favorited,
		},
		HTMLContent: post.HTMLContent,
	}

	comments, err := s.comments.ListByPost(postID)
	if err != nil {
		return nil, nil, err
	}
	views, err := s.commentViews(comments, viewerID)
	if err != nil {
		return nil, nil, err
	}
	return detail, views, nil
}

func (s *postService) commentViews(comments []model.Comment, viewerID string) ([]viewmodel.CommentView, error) {
	if len(comments) == 0 {
		return []viewmodel.CommentView{}, nil
	}

	userIDs := make([]string, 0, len(comments)*2)
	commentIDs := make([]int64, 0, len(comments))
	for _, comment := range comments {
		userIDs = append(userIDs, comment.AuthorID)
		if comment.ReplyToUserID != nil {
			userIDs = append(userIDs, *comment.ReplyToUserID)
		}
		commentIDs = append(commentIDs, comment.ID)
	}

	profiles, err := s.profiles.ListByIDs(userIDs)
	if err != nil {
		return nil, err
	}
	liked := map[int64]bool{}
	if viewerID != "" {
		if liked, err = s.likes.LikedComments(viewerID, commentIDs); err != nil {
			return nil, err
		}
	}

	roots := make([]viewmodel.CommentView, 0)
	rootIndex := make(map[int64]int)
	for _, comment := range comments {
		if comment.ParentID != nil {
			continue
		}
		view := commentView(comment, profiles, liked)
		rootIndex[comment.ID] = len(roots)
		roots = append(roots, view)
	}
	for _, comment := range comments {
		if comment.ParentID == nil {
			continue
		}
		index, ok := rootIndex[*comment.ParentID]
		if !ok {
			continue
		}
		roots[index].Replies = append(roots[index].Replies, commentView(comment, profiles, liked))
	}
	return roots, nil
}

func commentView(comment model.Comment, profiles map[string]model.UserProfile, liked map[int64]bool) viewmodel.CommentView {
	view := viewmodel.CommentView{
		ID:          comment.ID,
		PostID:      comment.PostID,
		ParentID:    comment.ParentID,
		Author:      viewmodel.NewUserSummary(comment.AuthorID, profilePointer(profiles, comment.AuthorID)),
		ContentHTML: comment.ContentHTML,
		ContentText: content.Excerpt(comment.ContentMarkdown, 200),
		CreatedAt:   comment.CreatedAt,
		LikeCount:   comment.LikeCount,
		Liked:       liked[comment.ID],
	}
	if comment.ReplyToUserID != nil {
		replyTo := viewmodel.NewUserSummary(*comment.ReplyToUserID, profilePointer(profiles, *comment.ReplyToUserID))
		view.ReplyTo = &replyTo
	}
	return view
}

func normalizePage(pageNo, pageSize int) (int, int) {
	if pageSize <= 0 || pageSize > 100 {
		pageSize = PageSize
	}
	if pageNo < 1 {
		pageNo = 1
	}
	return pageNo, pageSize
}

func TotalPages(total int64, pageSize int) int {
	if total <= 0 || pageSize <= 0 {
		return 1
	}
	pages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if pages < 1 {
		return 1
	}
	return pages
}
