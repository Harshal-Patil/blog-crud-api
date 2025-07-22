package repository

import (
	"blog-crud-api/models"

	"gorm.io/gorm"
)

type BlogRepository interface {
	Create(post *models.BlogPost) error
	GetAll() ([]models.BlogPost, error)
	GetByID(id uint) (*models.BlogPost, error)
	Update(post *models.BlogPost) error
	Delete(id uint) error
	GetAllWithPagination(limit int, offset int) ([]models.BlogPost, error)
	GetByUser(userID int) ([]models.BlogPost, error)
}

type blogRepo struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) BlogRepository {
	return &blogRepo{db}
}
func (r *blogRepo) Create(post *models.BlogPost) error {
	return r.db.Create(post).Error
}

func (r *blogRepo) GetAll() ([]models.BlogPost, error) {
	var posts []models.BlogPost
	err := r.db.Find(&posts).Error
	return posts, err
}

func (r *blogRepo) GetByID(id uint) (*models.BlogPost, error) {
	var post models.BlogPost
	err := r.db.First(&post, id).Error
	return &post, err
}

func (r *blogRepo) Update(post *models.BlogPost) error {
	return r.db.Save(post).Error
}

func (r *blogRepo) Delete(id uint) error {
	return r.db.Delete(&models.BlogPost{}, id).Error
}
func (r *blogRepo) GetAllWithPagination(limit int, offset int) ([]models.BlogPost, error) {
	var posts []models.BlogPost
	err := r.db.Limit(limit).Offset(offset).Order("created_at desc").Find(&posts).Error
	return posts, err
}

func (r *blogRepo) GetByUser(userID int) ([]models.BlogPost, error) {
	var posts []models.BlogPost
	err := r.db.Where("user_id = ?", userID).Find(&posts).Error
	return posts, err
}
