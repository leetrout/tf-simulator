package api

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/leetrout/terraform-sim/internal/models"
)

type SubnetsHandler struct {
	queries *models.Queries
	logger  *slog.Logger
}

func NewSubnetsHandler(queries *models.Queries, logger *slog.Logger) *SubnetsHandler {
	return &SubnetsHandler{
		queries: queries,
		logger:  logger,
	}
}

func (h *SubnetsHandler) RegisterRoutes(e *echo.Group) {
	e.GET("/subnets", h.listSubnets)
	e.GET("/subnets/:id", h.getSubnetByID)
	e.POST("/subnets", h.createSubnet)
}

type subnetResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	NetworkID int64  `json:"network_id"`
}

func toSubnetResponse(s models.Subnet) subnetResponse {
	return subnetResponse{
		ID:        s.ID,
		Name:      s.Name,
		NetworkID: s.NetworkID,
	}
}

type subnetsResponse []subnetResponse

func toSubnetsResponse(subnets []models.Subnet) subnetsResponse {
	resp := make(subnetsResponse, len(subnets))
	for i, s := range subnets {
		resp[i] = toSubnetResponse(s)
	}
	return resp
}

func (h *SubnetsHandler) listSubnets(c echo.Context) error {
	subnets, err := h.queries.ListSubnets(c.Request().Context())
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, toSubnetsResponse(subnets))
}

type getSubnetByIDRequest struct {
	ID int64 `param:"id"`
}

func (h *SubnetsHandler) getSubnetByID(c echo.Context) error {
	var req getSubnetByIDRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	subnet, err := h.queries.GetSubnet(c.Request().Context(), req.ID)
	if err == sql.ErrNoRows {
		return echo.ErrNotFound
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, toSubnetResponse(subnet))
}

type createSubnetRequest struct {
	Name      string `json:"name"`
	NetworkID int64  `json:"network_id"`
}

func (h *SubnetsHandler) createSubnet(c echo.Context) error {
	var req createSubnetRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	params := models.CreateSubnetParams{
		Name:      req.Name,
		NetworkID: req.NetworkID,
	}

	subnet, err := h.queries.CreateSubnet(c.Request().Context(), params)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, toSubnetResponse(subnet))
}
