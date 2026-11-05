package repository

import (
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model"

	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *model.Post) error
	GetById(id int64) (*model.Post, error)
	List(offset, limit int) ([]model.Post, error)
	Update(post *model.Post) error
	Delete(id int64) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) GetById(id int64) (*model.Post, error) {
	var post model.Post
	//if err := r.db.Preload("Comments").First(&post, id).Error; err != nil {
	//	return nil, err
	//}
	if err := r.db.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) List(offset int, limit int) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.Offset(offset).Limit(limit).
		Find(&posts).Error
	return posts, err
}

func (r *postRepository) Update(post *model.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id int64) error {
	return r.db.Delete(&model.Post{}, id).Error
}
