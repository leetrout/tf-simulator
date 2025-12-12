package api

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/leetrout/terraform-sim/internal/models"
)

type DatabasesHandler struct {
	queries *models.Queries
	logger  *slog.Logger
}

func NewDatabasesHandler(queries *models.Queries, logger *slog.Logger) *DatabasesHandler {
	return &DatabasesHandler{
		queries: queries,
		logger:  logger,
	}
}

func (h *DatabasesHandler) RegisterRoutes(e *echo.Group) {
	e.GET("/databases", h.listDatabases)
	e.GET("/databases/:id", h.getDatabaseByID)
	e.POST("/databases", h.createDatabase)
}

type databaseResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	SubnetID int64  `json:"subnet_id"`
}

func toDatabaseResponse(d models.Database) databaseResponse {
	return databaseResponse{
		ID:       d.ID,
		Name:     d.Name,
		SubnetID: d.SubnetID,
	}
}

type databasesResponse []databaseResponse

func toDatabasesResponse(databases []models.Database) databasesResponse {
	resp := make(databasesResponse, len(databases))
	for i, d := range databases {
		resp[i] = toDatabaseResponse(d)
	}
	return resp
}

func (h *DatabasesHandler) listDatabases(c echo.Context) error {
	databases, err := h.queries.ListDatabases(c.Request().Context())
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, toDatabasesResponse(databases))
}

type getDatabaseByIDRequest struct {
	ID int64 `param:"id"`
}

func (h *DatabasesHandler) getDatabaseByID(c echo.Context) error {
	var req getDatabaseByIDRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	database, err := h.queries.GetDatabase(c.Request().Context(), req.ID)
	if err == sql.ErrNoRows {
		return echo.ErrNotFound
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, toDatabaseResponse(database))
}

type createDatabaseRequest struct {
	Name     string `json:"name"`
	SubnetID int64  `json:"subnet_id"`
}

func (h *DatabasesHandler) createDatabase(c echo.Context) error {
	var req createDatabaseRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	params := models.CreateDatabaseParams{
		Name:     req.Name,
		SubnetID: req.SubnetID,
	}

	database, err := h.queries.CreateDatabase(c.Request().Context(), params)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, toDatabaseResponse(database))
}
