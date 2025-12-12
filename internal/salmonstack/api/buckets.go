package api

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/leetrout/terraform-sim/internal/models"
)

type BucketsHandler struct {
	queries *models.Queries
	logger  *slog.Logger
}

func NewBucketsHandler(queries *models.Queries, logger *slog.Logger) *BucketsHandler {
	return &BucketsHandler{
		queries: queries,
		logger:  logger,
	}
}

func (h *BucketsHandler) RegisterRoutes(e *echo.Group) {
	e.GET("/buckets", h.listBuckets)
	e.GET("/buckets/:id", h.getBucketByID)
	e.POST("/buckets", h.createBucket)
}

type bucketResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func toBucketResponse(b models.Bucket) bucketResponse {
	return bucketResponse{
		ID:   b.ID,
		Name: b.Name,
	}
}

type bucketsResponse []bucketResponse

func toBucketsResponse(buckets []models.Bucket) bucketsResponse {
	resp := make(bucketsResponse, len(buckets))
	for i, b := range buckets {
		resp[i] = toBucketResponse(b)
	}
	return resp
}

func (h *BucketsHandler) listBuckets(c echo.Context) error {
	buckets, err := h.queries.ListBuckets(c.Request().Context())
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, toBucketsResponse(buckets))
}

type getBucketByIDRequest struct {
	ID int64 `param:"id"`
}

func (h *BucketsHandler) getBucketByID(c echo.Context) error {
	var req getBucketByIDRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	bucket, err := h.queries.GetBucket(c.Request().Context(), req.ID)
	if err == sql.ErrNoRows {
		return echo.ErrNotFound
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, toBucketResponse(bucket))
}

type createBucketRequest struct {
	Name string `json:"name"`
}

func (h *BucketsHandler) createBucket(c echo.Context) error {
	var req createBucketRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	bucket, err := h.queries.CreateBucket(c.Request().Context(), req.Name)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, toBucketResponse(bucket))
}
