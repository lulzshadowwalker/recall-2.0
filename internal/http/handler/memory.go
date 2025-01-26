package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lulzshadowwalker/recall/internal"
)

type Memory struct {
	Repository MemoryRepository
}

//  NOTE: We might want to add a service later on.
type MemoryRepository interface {
	CreateMemory(c context.Context, params internal.CreateMemoryParams) (internal.Memory, error)
	DeleteMemory(c context.Context, id int) error
	GetMemories(c context.Context) ([]internal.Memory, error)
}

func NewMemoryHandler(r MemoryRepository) *Memory {
	return &Memory{r}
}

func (mh *Memory) Index(c echo.Context) error {
	m, err := mh.Repository.GetMemories(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, mh.collection(m))
}

type CreateMemoryRequest struct {
	Content string `json:"content" validate:"required"`
}

func (mh *Memory) Create(c echo.Context) error {
	var req CreateMemoryRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	mem, err := mh.Repository.CreateMemory(c.Request().Context(), internal.CreateMemoryParams{
		Content: req.Content,
    UserID: 2, //  TODO: Get user ID from auth
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, mh.resource(mem))
}

type DeleteMemoryRequest struct {
	ID int `param:"memory"`
}

func (mh *Memory) Delete(c echo.Context) error {
	var req DeleteMemoryRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	if err := mh.Repository.DeleteMemory(c.Request().Context(), req.ID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (mh *Memory) resource(m internal.Memory) echo.Map {
	return echo.Map{
		"id":   strconv.Itoa(m.ID),
		"type": "memory",
		"attributes": echo.Map{
			"content":   m.Content,
			"createdAt": m.CreatedAt,
			"updatedAt": m.UpdatedAt,
		},
		"relationships": echo.Map{},
	}
}

func (mh *Memory) collection(m []internal.Memory) echo.Map {
	res := make([]echo.Map, len(m))
	for i, mem := range m {
		res[i] = mh.resource(mem)
	}

	return echo.Map{
		"data": res,
	}
}
