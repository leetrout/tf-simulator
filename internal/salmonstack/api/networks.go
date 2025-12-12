package api

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/leetrout/terraform-sim/internal/models"
)

type NetworksHandler struct {
	queries *models.Queries
	logger  *slog.Logger
}

func NewNetworksHandler(queries *models.Queries, logger *slog.Logger) *NetworksHandler {
	return &NetworksHandler{
		queries: queries,
		logger:  logger,
	}
}

func (h *NetworksHandler) RegisterRoutes(e *echo.Group) {
	e.GET("/networks", h.listNetworks)
	e.GET("/networks/:id", h.getNetworkByID)
	e.POST("/networks", h.createNetwork)
}

type networkResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func toNetworkResponse(n models.Network) networkResponse {
	return networkResponse{
		ID:   n.ID,
		Name: n.Name,
	}
}

type networksResponse []networkResponse

func toNetworksResponse(networks []models.Network) networksResponse {
	resp := make(networksResponse, len(networks))
	for i, n := range networks {
		resp[i] = toNetworkResponse(n)
	}
	return resp
}

func (h *NetworksHandler) listNetworks(c echo.Context) error {
	networks, err := h.queries.ListNetworks(c.Request().Context())
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, toNetworksResponse(networks))
}

type getNetworkByIDRequest struct {
	ID int64 `param:"id"`
}

func (h *NetworksHandler) getNetworkByID(c echo.Context) error {
	var req getNetworkByIDRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	network, err := h.queries.GetNetwork(c.Request().Context(), req.ID)
	if err == sql.ErrNoRows {
		return echo.ErrNotFound
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, toNetworkResponse(network))
}

type createNetworkRequest struct {
	Name string `json:"name"`
}

func (h *NetworksHandler) createNetwork(c echo.Context) error {
	var req createNetworkRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	network, err := h.queries.CreateNetwork(c.Request().Context(), req.Name)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, toNetworkResponse(network))
}
