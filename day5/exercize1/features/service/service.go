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
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) GetAll(ctx context.Context) ([]core.Tasks, error) {
	logger, loggerClose, err := core.NewLogs("DEBUG")
	if err != nil {
		panic(err)
	}
	defer loggerClose()
	task, err := s.repo.GetAll(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("%sError:", err))
		return nil, err
	}
	if task == nil {
		logger.Info("nil Tasks respons user")
		return []core.Tasks{}, core.ErrNotFound
	}
	return task, nil
}

func (s *TaskService) GetByID(ctx context.Context, id int) (core.Tasks, error) {
	logger, loggerClose, err := core.NewLogs("DEBUG")
	if err != nil {
		panic(err)
	}
	defer loggerClose()
	if id <= 0 {
		logger.Info(" users bad request id ")
		return core.Tasks{}, core.ErrBadRequest
	}
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Error(fmt.Sprintf("%s Error:", err))
		return core.Tasks{}, err
	}
	return task, nil

}

func (s *TaskService) Create(ctx context.Context, t core.Tasks) (core.Tasks, error) {
	logger, loggerClose, err := core.NewLogs("DEBUG")
	if err != nil {
		panic(err)
	}
	defer loggerClose()
	t.Title = strings.TrimSpace(t.Title)
	t.Description = strings.TrimSpace(t.Description)

	if len(t.Title) == 0 {
		logger.Info("user bad request title")
		return core.Tasks{}, core.ErrBadRequest
	}
	task, err := s.repo.Created(ctx, t)
	if err != nil {
		logger.Error(fmt.Sprintf("%sError:", err))
		return core.Tasks{}, err
	}

	return task, nil
}
func (s *TaskService) Delete(ctx context.Context, id int) error {
	logger, loggerClose, err := core.NewLogs("DEBUG")
	if err != nil {
		panic(err)
	}
	defer loggerClose()
	if id <= 0 {
		logger.Info("User Bad Request")
		return core.ErrBadRequest
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Info("task 404 Not Found")
			return core.ErrNotFound
		}
		logger.Error(fmt.Sprintf("%sError:", err))
		return err
	}
	return nil
}

func (s *TaskService) UpdateStatus(ctx context.Context, completed bool, id int) (core.Tasks, error) {
	logger, loggerClose, err := core.NewLogs("DEBUG")
	if err != nil {
		panic(err)
	}
	defer loggerClose()
	if id <= 0 {
		logger.Info("Users Bad Request")
		return core.Tasks{}, core.ErrBadRequest
	}
	time := time.Now()
	task, err := s.repo.UpdateStatus(ctx, id, completed, &time)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Info("task not found")
			return core.Tasks{}, core.ErrNotFound
		}
		logger.Error(fmt.Sprintf("%sError:", err))
		return core.Tasks{}, err
	}
	return task, nil
}

func (s *TaskService) UpdateContent(ctx context.Context, id int, title string, description string) (core.Tasks, error) {
	logger, loggerClose, err := core.NewLogs("DEBUG")
	if err != nil {
		panic(err)
	}
	defer loggerClose()
	if id <= 0 {
		logger.Info("Users Bad Request")
		return core.Tasks{}, core.ErrBadRequest
	}
	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)
	if len(title) == 0 && len(description) == 0 {
		logger.Info("user bad request title")
		return core.Tasks{}, core.ErrBadRequest
	}
	task, err := s.repo.UpdateContent(ctx, title, id, description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Info("task not found")
			return core.Tasks{}, core.ErrNotFound
		}
		logger.Error(fmt.Sprintf("%sError:", err))
		return core.Tasks{}, err
	}
	return task, nil
}
