package usecase

import (
	"boilerplate/internal/model"
	"boilerplate/internal/repository"
)

type UserUsecase interface {
	GetUserData(id int) (*model.User, error)
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) GetUserData(id int) (*model.User, error) {
	return u.repo.FindByID(id)
}