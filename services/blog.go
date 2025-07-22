package services

import (
	"blog-crud-api/models"
	"blog-crud-api/repository"
)

type BlogService interface {
	Create(post models.BlogPost) (*models.BlogPost, error)
	GetAll() ([]models.BlogPost, error)
	GetByID(id uint) (*models.BlogPost, error)
	Update(id uint, data models.BlogPost) (*models.BlogPost, error)
	Delete(id uint) error
	GetAllWithPagination(limit int, offset int) ([]models.BlogPost, error)
	GetByUser(userID int) ([]models.BlogPost, error)
}

type blogService struct {
	repo repository.BlogRepository
}

func NewBlogService(r repository.BlogRepository) BlogService {
	return &blogService{repo: r}
}

func (s *blogService) Create(post models.BlogPost) (*models.BlogPost, error) {
	err := s.repo.Create(&post)
	return &post, err
}

func (s *blogService) GetAll() ([]models.BlogPost, error) {
	return s.repo.GetAll()
}

func (s *blogService) GetByID(id uint) (*models.BlogPost, error) {
	return s.repo.GetByID(id)
}
func (s *blogService) Update(id uint, data models.BlogPost) (*models.BlogPost, error) {
	post, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	post.Title = data.Title
	post.Description = data.Description
	post.Body = data.Body
	err = s.repo.Update(post)
	return post, err
}

func (s *blogService) Delete(id uint) error {
	return s.repo.Delete(id)
}
func (s *blogService) GetAllWithPagination(limit int, offset int) ([]models.BlogPost, error) {
	return s.repo.GetAllWithPagination(limit, offset)
}

func (s *blogService) GetByUser(userID int) ([]models.BlogPost, error) {
	return s.repo.GetByUser(userID)
}
