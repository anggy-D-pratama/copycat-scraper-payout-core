package service

import (
	"boilerplate/internal/model"
	"boilerplate/internal/repository"
)

type UserService interface {
	GetUserData(id int) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUserData(id int) (*model.User, error) {
	return s.repo.FindByID(id)
}