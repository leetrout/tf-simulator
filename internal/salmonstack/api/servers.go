package api

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/leetrout/terraform-sim/internal/models"
)

type ServersHandler struct {
	queries *models.Queries
	logger  *slog.Logger
}

func NewServersHandler(queries *models.Queries, logger *slog.Logger) *ServersHandler {
	return &ServersHandler{
		queries: queries,
		logger:  logger,
	}
}

func (h *ServersHandler) RegisterRoutes(e *echo.Group) {
	e.GET("/servers", h.listServers)
	e.GET("/servers/:id", h.getServerByID)
	e.POST("/servers", h.createServer)
}

type serverResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func toServerResponse(s models.Server) serverResponse {
	return serverResponse{
		ID:   s.ID,
		Name: s.Name,
	}
}

type serversResponse []serverResponse

func toServersResponse(servers []models.Server) serversResponse {
	resp := make(serversResponse, len(servers))
	for i, s := range servers {
		resp[i] = toServerResponse(s)
	}
	return resp
}

func (h *ServersHandler) listServers(c echo.Context) error {
	servers, err := h.queries.ListServers(c.Request().Context())
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, toServersResponse(servers))
}

type getServerByIDRequest struct {
	ID int64 `param:"id"`
}

func (h *ServersHandler) getServerByID(c echo.Context) error {
	var req getServerByIDRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	server, err := h.queries.GetServer(c.Request().Context(), req.ID)
	if err == sql.ErrNoRows {
		return echo.ErrNotFound
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, toServerResponse(server))
}

type createServerRequest struct {
	Name       string `json:"name"`
	SubnetID   int64  `json:"subnet_id"`
	StaticIPID *int64 `json:"static_ip_id,omitempty"`
}

func (h *ServersHandler) createServer(c echo.Context) error {
	var req createServerRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	var staticIPID sql.NullInt64
	if req.StaticIPID != nil {
		staticIPID = sql.NullInt64{
			Int64: *req.StaticIPID,
			Valid: true,
		}
	}

	params := models.CreateServerParams{
		Name:       req.Name,
		SubnetID:   req.SubnetID,
		StaticIpID: staticIPID,
	}

	server, err := h.queries.CreateServer(c.Request().Context(), params)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, toServerResponse(server))
}
