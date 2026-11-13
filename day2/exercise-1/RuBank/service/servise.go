package service

import (
	"context"
	repository "study/RuBank/Repository"
	"study/RuBank/models"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Transfer(ctx context.Context, fromID int, toID int, amount int) error {
	if amount <= 0 {
		return models.ErrSmallAmount
	}
	const num = 1000000
	if amount >= num {
		return models.ErrSuspicious
	}
	if err := s.repo.Transfer(ctx, fromID, toID, amount); err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetId(ctx context.Context, id int) (*models.Users, error) {
	NewUser, err := s.repo.GetId(ctx, id)
	if err != nil {
		return nil, err
	}
	return NewUser, nil
}

func (s *UserService) CreateAcc(ctx context.Context, u *models.Users) error {
	if err := s.repo.CreateAccount(ctx, u); err != nil {
		return err
	}
	return nil
}
