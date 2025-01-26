package psql


import (
	"context"
	"fmt"

	"github.com/lulzshadowwalker/recall/internal"
	"github.com/lulzshadowwalker/recall/internal/psql/db"
)

type Memory struct {
	q *db.Queries
}

func NewMemory(d db.DBTX) *Memory {
	return &Memory{
		q: db.New(d),
	}
}

func toEntity(m db.Memory) internal.Memory {
  return internal.Memory{
    ID:        int(m.ID),
    Content:   m.Content,
    UserID:    int(m.UserID),
    CreatedAt: m.CreatedAt.Time,
    UpdatedAt: m.UpdatedAt.Time,
  }
}

func (m *Memory) CreateMemory(c context.Context, params internal.CreateMemoryParams) (internal.Memory, error) {
  res, err := m.q.CreateMemory(c, db.CreateMemoryParams{
    Content: params.Content,
    UserID: int64(params.UserID),
  })

  if err != nil {
    return internal.Memory{}, fmt.Errorf("failed to create memory because %w", err)
  }

  return toEntity(res), nil
}

func (m *Memory) DeleteMemory(c context.Context, id int) error {
  if err :=  m.q.DeleteMemory(c, int64(id)); err != nil {
    return fmt.Errorf("failed to delete memory because %w", err)
  }

  return nil
}

func (m *Memory) GetMemories(c context.Context) ([]internal.Memory, error) {
	res, err := m.q.GetMemories(c)
	if err != nil {
		return nil, fmt.Errorf("failed to get memories from db because %w", err)
	}

	memories := make([]internal.Memory, len(res))
	for i, r := range res {
		memories[i] = toEntity(r)
	}

	return memories, nil
}
