package salmonstack

import (
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/leetrout/terraform-sim/internal/models"
	"github.com/leetrout/terraform-sim/internal/salmonstack/api"
	slogecho "github.com/samber/slog-echo"
)

func parseLogLevel(level string) (slog.Level, error) {
	var l slog.Level
	err := l.UnmarshalText([]byte(level))
	return l, err
}

func StartServer() error {
	logLevel, err := parseLogLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		println("Invalid LOG_LEVEL, defaulting to INFO")
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	db, err := openSQLite()
	if err != nil {
		return err
	}

	queries := models.New(db)

	e := echo.New()

	e.Use(middleware.RequestID())
	e.Use(slogecho.New(logger))
	e.Use(middleware.Recover())

	apiGroup := e.Group("/api/v1")

	networksHandler := api.NewNetworksHandler(queries, logger)
	networksHandler.RegisterRoutes(apiGroup)

	subnetsHandler := api.NewSubnetsHandler(queries, logger)
	subnetsHandler.RegisterRoutes(apiGroup)

	staticIPsHandler := api.NewStaticIPsHandler(queries, logger)
	staticIPsHandler.RegisterRoutes(apiGroup)

	domainNamesHandler := api.NewDomainNamesHandler(queries, logger)
	domainNamesHandler.RegisterRoutes(apiGroup)

	serversHandler := api.NewServersHandler(queries, logger)
	serversHandler.RegisterRoutes(apiGroup)

	loadBalancersHandler := api.NewLoadBalancersHandler(queries, logger)
	loadBalancersHandler.RegisterRoutes(apiGroup)

	databasesHandler := api.NewDatabasesHandler(queries, logger)
	databasesHandler.RegisterRoutes(apiGroup)

	bucketsHandler := api.NewBucketsHandler(queries, logger)
	bucketsHandler.RegisterRoutes(apiGroup)

	return e.Start(":8000")
}
