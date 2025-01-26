package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"time"

	"github.com/labstack/echo/v4"
	"github.com/lulzshadowwalker/recall/internal"
	"github.com/lulzshadowwalker/recall/internal/http/app"
	"github.com/lulzshadowwalker/recall/internal/http/handler"
)

type MockMemoryRepository struct {
	memories []internal.Memory
}

func NewMockMemoryRepository() *MockMemoryRepository {
	created, _ := time.Parse(time.RFC3339, "2025-01-24T20:54:51.431562+03:00")
	updated, _ := time.Parse(time.RFC3339, "2025-01-24T20:54:51.431562+03:00")

	return &MockMemoryRepository{
		memories: []internal.Memory{
			{
				ID:        0,
				Content:   "Hello, world!",
				CreatedAt: created,
				UpdatedAt: updated,
				UserID:    1,
			},
		},
	}
}

func (r *MockMemoryRepository) CreateMemory(ctx context.Context, params internal.CreateMemoryParams) (internal.Memory, error) {
	created, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	updated, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")

	m := internal.Memory{
		ID:        1,
		Content:   params.Content,
		CreatedAt: created,
		UpdatedAt: updated,
		UserID:    1,
	}

	r.memories = append(r.memories, m)

	return m, nil
}

func (r *MockMemoryRepository) DeleteMemory(ctx context.Context, id int) error {
	found := false
	for i, m := range r.memories {
		if m.ID == id {
			r.memories = append(r.memories[:i], r.memories[i+1:]...)
			found = true
			break
		}
	}

	if found {
		return nil
	}

	return errors.New("Memory not found")
}

func (r *MockMemoryRepository) GetMemories(ctx context.Context) ([]internal.Memory, error) {
	return r.memories, nil
}

func TestGetMemories(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/memories", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	mh := handler.NewMemoryHandler(NewMockMemoryRepository())
	err := mh.Index(ctx)

	body := rec.Body.String()
	t.Logf("%s", body)

	if err != nil {
		t.Errorf("error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", rec.Code)
	}

	var got, want interface{}
	json.Unmarshal([]byte(`{"data":[{"attributes":{"content":"Hello, world!","createdAt":"2025-01-24T20:54:51.431562+03:00","updatedAt":"2025-01-24T20:54:51.431562+03:00"},"id":"0","relationships":{},"type":"memory"}]}`), &want)
	json.Unmarshal([]byte(body), &got)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

//  TODO: Test validation
func TestCreateMemory(t *testing.T) {
  e := echo.New()
  e.Validator = app.NewRecallValidator()
  req := httptest.NewRequest(http.MethodPost, "/memories", strings.NewReader(`{"content":"Hello, lulzie!"}`))
  req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
  rec := httptest.NewRecorder()
  ctx := e.NewContext(req, rec)

  mh := handler.NewMemoryHandler(NewMockMemoryRepository())
  err := mh.Create(ctx)

  body := rec.Body.String()
  t.Logf("%s", body)

  if err != nil {
    t.Errorf("error: %v", err)
  }

  if rec.Code != http.StatusCreated {
    t.Errorf("expected status code 201, got %d", rec.Code)
  }

  var got, want interface{}
  json.Unmarshal([]byte(`{"attributes":{"content":"Hello, lulzie!","createdAt":"2021-01-01T00:00:00Z","updatedAt":"2021-01-01T00:00:00Z"},"id":"1","relationships":{},"type":"memory"}`), &want)
  json.Unmarshal([]byte(body), &got)

  if !reflect.DeepEqual(got, want) {
    t.Errorf("expected %v, got %v", want, got)
  }
}

func TestDeleteMemory(t *testing.T) {
  e := echo.New()
  e.Validator = app.NewRecallValidator()
  req := httptest.NewRequest(http.MethodDelete, "/memories/0", nil)
  rec := httptest.NewRecorder()
  ctx := e.NewContext(req, rec)

  mh := handler.NewMemoryHandler(NewMockMemoryRepository())
  err := mh.Delete(ctx)

  body := rec.Body.String()
  t.Logf("%s", body)

  if err != nil {
    t.Errorf("error: %v", err)
  }

  if rec.Code != http.StatusOK {
    t.Errorf("expected status code 200, got %d", rec.Code)
  }

  if body != "" {
    t.Errorf("expected empty body, got %s", body)
  }
}
