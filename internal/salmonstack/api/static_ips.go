package api

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/leetrout/terraform-sim/internal/models"
)

type StaticIPsHandler struct {
	queries *models.Queries
	logger  *slog.Logger
}

func NewStaticIPsHandler(queries *models.Queries, logger *slog.Logger) *StaticIPsHandler {
	return &StaticIPsHandler{
		queries: queries,
		logger:  logger,
	}
}

func (h *StaticIPsHandler) RegisterRoutes(e *echo.Group) {
	e.GET("/static-ips", h.listStaticIPs)
	e.GET("/static-ips/:id", h.getStaticIPByID)
	e.POST("/static-ips", h.createStaticIP)
}

type staticIPResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
}

func toStaticIPResponse(s models.StaticIp) staticIPResponse {
	return staticIPResponse{
		ID:        s.ID,
		Name:      s.Name,
		IPAddress: s.IpAddress,
	}
}

type staticIPsResponse []staticIPResponse

func toStaticIPsResponse(staticIPs []models.StaticIp) staticIPsResponse {
	resp := make(staticIPsResponse, len(staticIPs))
	for i, s := range staticIPs {
		resp[i] = toStaticIPResponse(s)
	}
	return resp
}

func (h *StaticIPsHandler) listStaticIPs(c echo.Context) error {
	staticIPs, err := h.queries.ListStaticIPs(c.Request().Context())
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, toStaticIPsResponse(staticIPs))
}

type getStaticIPByIDRequest struct {
	ID int64 `param:"id"`
}

func (h *StaticIPsHandler) getStaticIPByID(c echo.Context) error {
	var req getStaticIPByIDRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	staticIP, err := h.queries.GetStaticIP(c.Request().Context(), req.ID)
	if err == sql.ErrNoRows {
		return echo.ErrNotFound
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, toStaticIPResponse(staticIP))
}

type createStaticIPRequest struct {
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
}

func (h *StaticIPsHandler) createStaticIP(c echo.Context) error {
	var req createStaticIPRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	params := models.CreateStaticIPParams{
		Name:      req.Name,
		IpAddress: req.IPAddress,
	}

	staticIP, err := h.queries.CreateStaticIP(c.Request().Context(), params)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, toStaticIPResponse(staticIP))
}
