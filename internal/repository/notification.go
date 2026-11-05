package repository

import (
	"gorm.io/gorm"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model"
)

type NotificationRepository interface {
	WithDB(db *gorm.DB) NotificationRepository
	CreateBatch(notifications []model.Notification) error
	ListByRecipient(recipientID string, offset, limit int) ([]model.Notification, error)
	CountByRecipient(recipientID string) (int64, error)
	CountUnread(recipientID string) (int64, error)
	MarkAllRead(recipientID string) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) WithDB(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) CreateBatch(notifications []model.Notification) error {
	if len(notifications) == 0 {
		return nil
	}
	return r.db.Create(&notifications).Error
}

func (r *notificationRepository) ListByRecipient(recipientID string, offset, limit int) ([]model.Notification, error) {
	var notifications []model.Notification
	err := r.db.Where("recipient_id = ?", recipientID).
		Order("created_at DESC, id DESC").
		Offset(offset).Limit(limit).
		Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) CountByRecipient(recipientID string) (int64, error) {
	var total int64
	err := r.db.Model(&model.Notification{}).Where("recipient_id = ?", recipientID).Count(&total).Error
	return total, err
}

func (r *notificationRepository) CountUnread(recipientID string) (int64, error) {
	var total int64
	err := r.db.Model(&model.Notification{}).
		Where("recipient_id = ? AND read_at IS NULL", recipientID).
		Count(&total).Error
	return total, err
}

func (r *notificationRepository) MarkAllRead(recipientID string) error {
	return r.db.Model(&model.Notification{}).
		Where("recipient_id = ? AND read_at IS NULL", recipientID).
		UpdateColumn("read_at", gorm.Expr("now()")).Error
}
