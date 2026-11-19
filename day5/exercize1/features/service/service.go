package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"study/core"
	"study/features/repository"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type TaskService struct {
	repo *repository.TaskRepository
	zap  *zap.Logger
}

func NewTaskService(repo *repository.TaskRepository, zap *zap.Logger) *TaskService {
	return &TaskService{repo: repo, zap: zap}
}

func (s *TaskService) GetAll(ctx context.Context) ([]core.Tasks, error) {
	task, err := s.repo.GetAll(ctx)
	if err != nil {
		s.zap.Error(fmt.Sprintf("%sError:", err))
		return nil, err
	}
	if task == nil {
		s.zap.Info("nil Tasks respons user")
		return []core.Tasks{}, nil
	}
	return task, nil
}

func (s *TaskService) GetByID(ctx context.Context, id int) (core.Tasks, error) {
	if id <= 0 {
		s.zap.Info(" users bad request id ")
		return core.Tasks{}, core.ErrBadRequest
	}
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.zap.Error(fmt.Sprintf("%s Error:", err))
		return core.Tasks{}, err
	}
	return task, nil

}

func (s *TaskService) Create(ctx context.Context, t core.Tasks) (core.Tasks, error) {
	t.Title = strings.TrimSpace(t.Title)
	t.Description = strings.TrimSpace(t.Description)

	if len(t.Title) == 0 {
		s.zap.Info("user bad request title")
		return core.Tasks{}, core.ErrBadRequest
	}
	task, err := s.repo.Created(ctx, t)
	if err != nil {
		s.zap.Error(fmt.Sprintf("%sError:", err))
		return core.Tasks{}, err
	}

	return task, nil
}
func (s *TaskService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		s.zap.Info("User Bad Request")
		return core.ErrBadRequest
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.zap.Info("task 404 Not Found")
			return core.ErrNotFound
		}
		s.zap.Error(fmt.Sprintf("%sError:", err))
		return err
	}
	return nil
}

func (s *TaskService) UpdateStatus(ctx context.Context, completed bool, id int) (core.Tasks, error) {
	if id <= 0 {
		s.zap.Info("Users Bad Request")
		return core.Tasks{}, core.ErrBadRequest
	}
	time := time.Now()
	task, err := s.repo.UpdateStatus(ctx, id, completed, &time)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.zap.Info("task not found")
			return core.Tasks{}, core.ErrNotFound
		}
		s.zap.Error(fmt.Sprintf("%sError:", err))
		return core.Tasks{}, err
	}
	return task, nil
}

func (s *TaskService) UpdateContent(ctx context.Context, id int, title string, description string) (core.Tasks, error) {
	if id <= 0 {
		s.zap.Info("Users Bad Request")
		return core.Tasks{}, core.ErrBadRequest
	}
	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)
	if title == "" {
		s.zap.Info("user bad request title")
		return core.Tasks{}, core.ErrBadRequest
	}
	task, err := s.repo.UpdateContent(ctx, title, id, description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.zap.Info("task not found")
			return core.Tasks{}, core.ErrNotFound
		}
		s.zap.Error(fmt.Sprintf("%sError:", err))
		return core.Tasks{}, err
	}
	return task, nil
}
