package api

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/leetrout/terraform-sim/internal/models"
)

type DomainNamesHandler struct {
	queries *models.Queries
	logger  *slog.Logger
}

func NewDomainNamesHandler(queries *models.Queries, logger *slog.Logger) *DomainNamesHandler {
	return &DomainNamesHandler{
		queries: queries,
		logger:  logger,
	}
}

func (h *DomainNamesHandler) RegisterRoutes(e *echo.Group) {
	e.GET("/domain-names", h.listDomainNames)
	e.GET("/domain-names/:id", h.getDomainNameByID)
	e.POST("/domain-names", h.createDomainName)
}

type domainNameResponse struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	StaticIPID *int64 `json:"static_ip_id,omitempty"`
}

func toDomainNameResponse(d models.DomainName) domainNameResponse {
	resp := domainNameResponse{
		ID:   d.ID,
		Name: d.Name,
	}
	if d.StaticIpID.Valid {
		resp.StaticIPID = &d.StaticIpID.Int64
	}
	return resp
}

type domainNamesResponse []domainNameResponse

func toDomainNamesResponse(domainNames []models.DomainName) domainNamesResponse {
	resp := make(domainNamesResponse, len(domainNames))
	for i, d := range domainNames {
		resp[i] = toDomainNameResponse(d)
	}
	return resp
}

func (h *DomainNamesHandler) listDomainNames(c echo.Context) error {
	domainNames, err := h.queries.ListDomainNames(c.Request().Context())
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, toDomainNamesResponse(domainNames))
}

type getDomainNameByIDRequest struct {
	ID int64 `param:"id"`
}

func (h *DomainNamesHandler) getDomainNameByID(c echo.Context) error {
	var req getDomainNameByIDRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	domainName, err := h.queries.GetDomainName(c.Request().Context(), req.ID)
	if err == sql.ErrNoRows {
		return echo.ErrNotFound
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, toDomainNameResponse(domainName))
}

type createDomainNameRequest struct {
	Name       string `json:"name"`
	StaticIPID *int64 `json:"static_ip_id,omitempty"`
}

func (h *DomainNamesHandler) createDomainName(c echo.Context) error {
	var req createDomainNameRequest
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

	params := models.CreateDomainNameParams{
		Name:       req.Name,
		StaticIpID: staticIPID,
	}

	domainName, err := h.queries.CreateDomainName(c.Request().Context(), params)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, toDomainNameResponse(domainName))
}
