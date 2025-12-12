package api

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/leetrout/terraform-sim/internal/models"
)

type LoadBalancersHandler struct {
	queries *models.Queries
	logger  *slog.Logger
}

func NewLoadBalancersHandler(queries *models.Queries, logger *slog.Logger) *LoadBalancersHandler {
	return &LoadBalancersHandler{
		queries: queries,
		logger:  logger,
	}
}

func (h *LoadBalancersHandler) RegisterRoutes(e *echo.Group) {
	e.GET("/load-balancers", h.listLoadBalancers)
	e.GET("/load-balancers/:id", h.getLoadBalancerByID)
	e.POST("/load-balancers", h.createLoadBalancer)
}

type loadBalancerResponse struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	SubnetID   int64  `json:"subnet_id"`
	StaticIPID *int64 `json:"static_ip_id,omitempty"`
}

func toLoadBalancerResponse(lb models.LoadBalancer) loadBalancerResponse {
	resp := loadBalancerResponse{
		ID:       lb.ID,
		Name:     lb.Name,
		SubnetID: lb.SubnetID,
	}
	if lb.StaticIpID.Valid {
		resp.StaticIPID = &lb.StaticIpID.Int64
	}
	return resp
}

type loadBalancersResponse []loadBalancerResponse

func toLoadBalancersResponse(loadBalancers []models.LoadBalancer) loadBalancersResponse {
	resp := make(loadBalancersResponse, len(loadBalancers))
	for i, lb := range loadBalancers {
		resp[i] = toLoadBalancerResponse(lb)
	}
	return resp
}

func (h *LoadBalancersHandler) listLoadBalancers(c echo.Context) error {
	loadBalancers, err := h.queries.ListLoadBalancers(c.Request().Context())
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, toLoadBalancersResponse(loadBalancers))
}

type getLoadBalancerByIDRequest struct {
	ID int64 `param:"id"`
}

func (h *LoadBalancersHandler) getLoadBalancerByID(c echo.Context) error {
	var req getLoadBalancerByIDRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	loadBalancer, err := h.queries.GetLoadBalancer(c.Request().Context(), req.ID)
	if err == sql.ErrNoRows {
		return echo.ErrNotFound
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, toLoadBalancerResponse(loadBalancer))
}

type createLoadBalancerRequest struct {
	Name       string `json:"name"`
	SubnetID   int64  `json:"subnet_id"`
	StaticIPID *int64 `json:"static_ip_id,omitempty"`
}

func (h *LoadBalancersHandler) createLoadBalancer(c echo.Context) error {
	var req createLoadBalancerRequest
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

	params := models.CreateLoadBalancerParams{
		Name:       req.Name,
		SubnetID:   req.SubnetID,
		StaticIpID: staticIPID,
	}

	loadBalancer, err := h.queries.CreateLoadBalancer(c.Request().Context(), params)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, toLoadBalancerResponse(loadBalancer))
}
