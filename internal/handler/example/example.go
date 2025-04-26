package example

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/dyleme/template/internal/domain"
)

type Service interface {
	Get(ctx context.Context, id int) (domain.Example, error)
	Update(ctx context.Context, params domain.Example) error
}
type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

func (h *Handler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	example, err := h.service.Get(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, viewExample(example))
}

func (h *Handler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id")
	}

	var params domain.Example
	if err := c.Bind(&params); err != nil {
		return err
	}

	params.ID = id

	err = h.service.Update(c.Request().Context(), params)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

type exampleView struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func viewExample(exmpl domain.Example) exampleView {
	return exampleView{
		ID:   exmpl.ID,
		Name: exmpl.Name,
	}
}
