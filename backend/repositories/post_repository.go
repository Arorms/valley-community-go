package repositories

import (
	"valley-community-go/models"

	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *models.Post) error
	GetById(id uint) (*models.Post, error)
	List(offset, limit int) ([]models.Post, error)
	Update(post *models.Post) error
	Delete(id uint) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) GetById(id uint) (*models.Post, error) {
	var post models.Post
	//if err := r.db.Preload("Comments").First(&post, id).Error; err != nil {
	//	return nil, err
	//}
	if err := r.db.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) List(offset int, limit int) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Offset(offset).Limit(limit).
		Find(&posts).Error
	return posts, err
}

func (r *postRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}
