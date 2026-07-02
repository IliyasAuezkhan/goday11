package service
import (
	"errors"
	"fmt"
	"day16/internal/domain"
)

type UserRepository interface {
	Save(user *domain.User) error
	FindByID(id int64) (*domain.User, error)
	Update(id int64, fields *domain.UpdateUserFields) error
	Delete(id int64) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(email, password string) (*domain.User, error) {
	if email == "" || password == "" {
		return nil, errors.New("invalid email or password")
	}

	user := &domain.User{Email: email, Password: password}
	if err := s.repo.Save(user); err != nil {
		return nil, fmt.Errorf("service failed: %w", err)
	}

	return user, nil
}

func (s *UserService) GetByID(id int64) (*domain.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}
	return s.repo.FindByID(id)
}

func (s *UserService) Update(id int64, email, password *string) error {
	if id <= 0 {
		return errors.New("invalid user ID")
	}

	fields := &domain.UpdateUserFields{
		Email:    email,
		Password: password,
	}

	return s.repo.Update(id, fields)
}

func (s *UserService) Delete(id int64) error {
	if id <= 0 {
		return errors.New("invalid user ID")
	}
	return s.repo.Delete(id)
}