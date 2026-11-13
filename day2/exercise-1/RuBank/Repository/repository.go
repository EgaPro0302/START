package repository

import (
	"context"
	"errors"
	"study/RuBank/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// создали пользователя
func (r *UserRepository) CreateAccount(ctx context.Context, u *models.Users) error {
	err := r.pool.QueryRow(ctx,
		"INSERT INTO users (id, money) VALUES ($1, $2) RETURNING id, money",
		u.ID, u.Money).Scan(&u.ID, &u.Money)
	return err
}

// можем посмтореть пользоватлея!
func (r *UserRepository) GetId(ctx context.Context, id int) (*models.Users, error) {
	var u models.Users
	rows := r.pool.QueryRow(ctx,
		"SELECT id, money FROM users WHERE id = $1",
		id,
	)
	if err := rows.Scan(&u.ID, &u.Money); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrUsersNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Transfer(ctx context.Context, fromID int, toID int, amount int) error {
	//делаем проверку
	if fromID == toID {
		return models.ErrYourself
	}
	if amount <= 0 {
		return models.ErrSmallAmount
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	tag1, err := tx.Exec(ctx,
		"UPDATE users SET money = money - $1 WHERE id = $2 AND money >= $1",
		amount, fromID,
	)
	if err != nil {
		return err
	}
	if tag1.RowsAffected() == 0 {
		return models.ErrInvalidTransfer

	}

	tag2, err := tx.Exec(ctx,
		"UPDATE users SET money = money + $1 WHERE id = $2 ",
		amount, toID,
	)
	if err != nil {
		return err
	}
	if tag2.RowsAffected() == 0 {
		return models.ErrInvalidTransfer
	}

	return tx.Commit(ctx)
}
