package cronx

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rizalgowandy/cronx/page"
	gdkMiddleware "github.com/rizalgowandy/gdk/pkg/httpx/echo/middleware"
)

const (
	QueryParamSort = "sort"
)

// SleepDuration defines the duration to sleep the server if the defined address is busy.
const SleepDuration = time.Second * 10

// NewServer creates a new HTTP server.
// - /			=> current server status.
// - /jobs		=> current jobs as frontend html.
// - /api/jobs	=> current jobs as json.
func NewServer(manager *Manager, address string) (*http.Server, error) {
	// Create server.
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(gdkMiddleware.RequestID())

	// Create server controller.
	ctrl := &ServerController{Manager: manager}

	// Register routes.
	e.GET("/", ctrl.HealthCheck)
	e.GET("/jobs", ctrl.Jobs)
	e.GET("/histories", ctrl.Histories)
	e.GET("/api/jobs", ctrl.APIJobs)
	e.GET("/api/histories", ctrl.APIHistories)

	return &http.Server{
		Addr:              address,
		Handler:           e,
		ReadHeaderTimeout: 60 * time.Second,
	}, nil
}

// NewSideCarServer creates a new sidecar HTTP server.
// HTTP server will be start automatically.
// - /			=> current server status.
// - /jobs		=> current jobs as frontend html.
// - /api/jobs	=> current jobs as json.
func NewSideCarServer(manager *Manager, address string) {
	// Create server.
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(gdkMiddleware.RequestID())

	// Create server controller.
	ctrl := &ServerController{Manager: manager}

	// Register routes.
	e.GET("/", ctrl.HealthCheck)
	e.GET("/jobs", ctrl.Jobs)
	e.GET("/histories", ctrl.Histories)
	e.GET("/api/jobs", ctrl.APIJobs)
	e.GET("/api/histories", ctrl.APIHistories)

	// Overcome issue with socket-master respawning 2nd app,
	// We will keep trying to run the server.
	// If the current address is busy,
	// sleep then try again until the address has become available.
	for {
		if err := e.Start(address); err != nil {
			time.Sleep(SleepDuration)
		}
	}
}

// ServerController is http server controller.
type ServerController struct {
	// Manager controls all the underlying job.
	Manager *Manager
}

// HealthCheck returns server status.
func (c *ServerController) HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, c.Manager.GetInfo())
}

// Jobs return job status as frontend template.
func (c *ServerController) Jobs(ctx echo.Context) error {
	index, err := page.GetJobsPageTemplate()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return index.Execute(
		ctx.Response().Writer,
		c.Manager.GetStatusData(ctx.QueryParam(QueryParamSort)),
	)
}

// APIJobs returns job status as json.
func (c *ServerController) APIJobs(ctx echo.Context) error {
	return ctx.JSON(
		http.StatusOK,
		c.Manager.GetStatusData(ctx.QueryParam(QueryParamSort)),
	)
}

// Histories return job history as frontend template.
func (c *ServerController) Histories(ctx echo.Context) error {
	index, err := page.GetHistoryTemplate()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	var req Request
	err = ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	req.url = *ctx.Request().URL

	data, err := c.Manager.GetHistoryData(ctx.Request().Context(), &req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return index.Execute(ctx.Response().Writer, data)
}

// APIHistories returns run histories as json.
func (c *ServerController) APIHistories(ctx echo.Context) error {
	var req Request
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	req.url = *ctx.Request().URL

	data, err := c.Manager.GetHistoryData(ctx.Request().Context(), &req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, data)
}
