package internal

import "time"

type Memory struct {
	ID        int
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    int
}

type CreateMemoryParams struct {
	Content string
	UserID  int
}
