package service

import (
	"user-backend/graph/model"
	"user-backend/repository"
)

type UserService struct {
	repository *repository.MongoRepository
}

func NewCustomerService(repository *repository.MongoRepository) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) CreateUser(user *model.User) (*model.User, error) {
	return s.repository.CreateUser(user)
}

func (s *UserService) UpdateUser(user *model.User) (*model.User, error) {
	return s.repository.UpdateUser(user)
}

func (s *UserService) DeleteUser(id string) (*model.Response, error) {
	err := s.repository.DeleteUser(id)
	if err != nil {
		return nil, err
	}

	return &model.Response{
		Message: "User deleted successfully",
		Status:  true,
	}, nil
}

func (s *UserService) GetUsers() ([]*model.User, error) {
	return s.repository.GetUsers()
}

func (s *UserService) GetUser(id string) (*model.User, error) {
	return s.repository.GetUser(id)
}
