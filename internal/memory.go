package internal

import "time"

type Memory struct {
	ID        int
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateMemoryParams struct {
	content string
}
