package repository

import (
	"context"
	"errors"
	"fmt"
	"study/core"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type TaskRepository struct {
	pool *pgxpool.Pool
	zap  *zap.Logger
}

func NewTaskRepository(pool *pgxpool.Pool, zap *zap.Logger) *TaskRepository {
	return &TaskRepository{pool: pool, zap: zap}
}

func (r *TaskRepository) GetAll(ctx context.Context) ([]core.Tasks, error) {

	rows, err := r.pool.Query(
		ctx,
		"SELECT id, title, description, completed, created_at, completed_at FROM tasks",
	)
	if err != nil {
		r.zap.Error(fmt.Sprintf("%sError:", err))
		return nil, err
	}
	defer rows.Close()
	var tasks []core.Tasks
	for rows.Next() {
		var t core.Tasks
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Completed,
			&t.CreatedAt,
			&t.CompletedAt,
		); err != nil {
			r.zap.Error(fmt.Sprintf("%sError:", err))
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	r.zap.Info("Get All Task")
	return tasks, nil
}

// показать task по id
func (r *TaskRepository) GetByID(ctx context.Context, id int) (core.Tasks, error) {
	var t core.Tasks
	if err := r.pool.QueryRow(ctx,
		"SELECT id, title, description, completed, created_at, completed_at FROM tasks WHERE id = $1",
		id,
	).Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Completed,
		&t.CreatedAt,
		&t.CompletedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) { // проверяем, что это именно "не найдено"
			return core.Tasks{}, core.ErrNotFound
		}
		r.zap.Error(fmt.Sprintf("%s   <-Error", err))
		return core.Tasks{}, err
	}
	r.zap.Info("GetById Task!")
	return t, nil
}

// создаем TASK
func (r *TaskRepository) Created(ctx context.Context, t core.Tasks) (core.Tasks, error) {
	if err := r.pool.QueryRow(
		ctx,
		`
		INSERT INTO tasks (title, description) 
		VALUES ($1, $2) 
		RETURNING id, title, description, completed, created_at, completed_at
		`,
		t.Title, t.Description,
	).Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Completed,
		&t.CreatedAt,
		&t.CompletedAt,
	); err != nil {
		r.zap.Error(fmt.Sprintf("%serror:", err))
		return core.Tasks{}, err
	}
	r.zap.Info("Created Task")

	return t, nil
}

// delete task
func (r *TaskRepository) Delete(ctx context.Context, id int) error {
	res, err := r.pool.Exec(
		ctx,
		"DELETE FROM tasks WHERE id =$1",
		id,
	)
	if err != nil {
		r.zap.Error(fmt.Sprintf("%sError:", err))
		return err
	}
	if res.RowsAffected() == 0 {
		r.zap.Info("Пользователь не нашел такую задачу!")
		return pgx.ErrNoRows
	}
	r.zap.Info("user Success delete task")
	return nil
}

// изменяем задачу на выполненую
func (r *TaskRepository) UpdateStatus(ctx context.Context, id int, completed bool, completed_at *time.Time) (core.Tasks, error) {
	var updateTask core.Tasks
	if err := r.pool.QueryRow(
		ctx,
		"UPDATE tasks SET completed = $1, completed_at = $2 WHERE id = $3 RETURNING id, title, description, completed, created_at, completed_at",
		completed, completed_at, id,
	).Scan(
		&updateTask.ID,
		&updateTask.Title,
		&updateTask.Description,
		&updateTask.Completed,
		&updateTask.CreatedAt,
		&updateTask.CompletedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.zap.Info("tast Not Found")
			return core.Tasks{}, pgx.ErrNoRows
		}
		r.zap.Error(fmt.Sprintf("%sError:", err))
		return core.Tasks{}, err
	}
	r.zap.Info("users Success completed task")
	return updateTask, nil
}

func (r *TaskRepository) UpdateContent(ctx context.Context, title string, id int, description string) (core.Tasks, error) {
	var updateTask core.Tasks
	if err := r.pool.QueryRow(
		ctx,
		"UPDATE tasks SET title =$1, description =$2 WHERE id = $3 RETURNING id, title, description, completed, created_at, completed_at",
		title, description, id,
	).Scan(
		&updateTask.ID,
		&updateTask.Title,
		&updateTask.Description,
		&updateTask.Completed,
		&updateTask.CreatedAt,
		&updateTask.CompletedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.zap.Info("tast Not Found")
			return core.Tasks{}, pgx.ErrNoRows
		}
		r.zap.Error(fmt.Sprintf("%sError", err))
		return core.Tasks{}, err
	}

	r.zap.Info("user Success Update Task")
	return updateTask, nil
}
