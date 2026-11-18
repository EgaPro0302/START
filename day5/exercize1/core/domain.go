package core

import "time"

type Tasks struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type DTOTask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
