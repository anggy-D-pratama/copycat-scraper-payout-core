package repository

import "boilerplate/internal/model"

type UserRepository interface {
	FindByID(id int) (*model.User, error)
}

type userRepository struct {
	// ponytail: ganti db connection misal db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) FindByID(id int) (*model.User, error) {
	// Mock return for now
	return &model.User{ID: id, Name: "John Doe", Email: "john@example.com"}, nil
}