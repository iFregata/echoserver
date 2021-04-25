package echox

import (
	"context"
	"echoserver/gorux"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type ServerConfig struct {
	Name         string `json:"name"`
	Developer    string `json:"developer"`
	Version      string `json:"version"`
	Branch       string `json:"branch"`
	Port         int    `json:"port"`
	ContextPath  string `json:"context_path"`
	LogLevel     string `json:"log_level"`
	EnableLogger bool   `json:"enable_logger"`
}

type EchoNi struct {
	*echo.Echo
}

type RoutingContext struct {
	echo.Context
}

type respBodyWrapper struct {
	StatusCode int         `json:"sc"`
	StatusText string      `json:"st"`
	Payload    interface{} `json:"payload,omitempty"`
}

var serverConfig *ServerConfig

func init() {
	serverConfig = loadServerConfig()
}

func New() EchoNi {
	echo := EchoNi{echo.New()}
	echo.Static(pathOf("/assets"), "assets")
	echo.HideBanner = true
	// Custom log level and setup Logger middleware logs
	// the information about each HTTP request.
	echo.customLog()
	// Setup the `liveness, readiness, inspect` routing
	echo.builtinRouting()
	return echo
}

func (e *EchoNi) GET(path string, fn func(c RoutingContext) error) {
	e.Echo.GET(pathOf(path), func(c echo.Context) error { return fn(RoutingContext{c}) })
}
func (e *EchoNi) POST(path string, fn func(c RoutingContext) error) {
	e.Echo.POST(pathOf(path), func(c echo.Context) error { return fn(RoutingContext{c}) })
}
func (e *EchoNi) PUT(path string, fn func(c RoutingContext) error) {
	e.Echo.PUT(pathOf(path), func(c echo.Context) error { return fn(RoutingContext{c}) })
}
func (e *EchoNi) PATCH(path string, fn func(c RoutingContext) error) {
	e.Echo.PATCH(pathOf(path), func(c echo.Context) error { return fn(RoutingContext{c}) })
}
func (e *EchoNi) DELETE(path string, fn func(c RoutingContext) error) {
	e.Echo.DELETE(pathOf(path), func(c echo.Context) error { return fn(RoutingContext{c}) })
}
func (e *EchoNi) HEAD(path string, fn func(c RoutingContext) error) {
	e.Echo.HEAD(pathOf(path), func(c echo.Context) error { return fn(RoutingContext{c}) })
}

func (e *EchoNi) Serve() {
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", serverConfig.Port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server")
		}
	}()

	// Make a chan to receive the sys signal to shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block the main goroutine waiting the shutdown signal
	<-quit
	e.Logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	e.Logger.Info("Server stutdown.")
}

func (rc RoutingContext) JsonWrap(payload interface{}, err error) error {
	if err != nil {
		rc.Logger().Error(err)
		return rc.JSON(200, &respBodyWrapper{StatusCode: 500, StatusText: err.Error()})
	}
	if payload == nil {
		return rc.JSON(200, &respBodyWrapper{StatusCode: 200, StatusText: "Ok"})
	}
	return rc.JSON(200, &respBodyWrapper{StatusCode: 200, StatusText: "Ok", Payload: &payload})
}

func (e *EchoNi) builtinRouting() {
	okFn := func(rc RoutingContext) error { return rc.Ok() }
	e.GET("/readiness", okFn)
	e.GET("/liveness", okFn)
	e.GET("/inspect", inspect)
}

// HTTP Response 200 OK
func (rc RoutingContext) Ok() error {
	return rc.NoContent(http.StatusOK)
}

// HTTP Response 201 Created
func (rc RoutingContext) Created() error {
	return rc.NoContent(http.StatusCreated)
}

// HTTP Response 202 Accepted
func (rc RoutingContext) Accepted() error {
	return rc.NoContent(http.StatusAccepted)
}

// HTTP Response 400 Bad Request
func (rc RoutingContext) BadRequest() error {
	return rc.NoContent(http.StatusBadRequest)
}

// HTTP Response 401 Unauthorized
func (rc RoutingContext) Unauthorized() error {
	return rc.NoContent(http.StatusUnauthorized)
}

// HTTP Response 403 Forbidden
func (rc RoutingContext) Forbidden() error {
	return rc.NoContent(http.StatusForbidden)
}

// HTTP Response 404 Not Found
func (rc RoutingContext) NotFound() error {
	return rc.NoContent(http.StatusNotFound)
}

// HTTP Response 500 Internal Server Error
func (rc RoutingContext) InternalServerError() error {
	return rc.NoContent(http.StatusInternalServerError)
}

// HTTP Response 501 Not Implemented
func (rc RoutingContext) NotImplemented() error {
	return rc.NoContent(http.StatusNotImplemented)
}

// HTTP Response 503 Service Unavailable
func (rc RoutingContext) ServiceUnavailable() error {
	return rc.NoContent(http.StatusServiceUnavailable)
}

// HTTP Response 504 Gateway Timeout
func (rc RoutingContext) Gatewayimeout() error {
	return rc.NoContent(http.StatusGatewayTimeout)
}

// Inspect server configuration
func inspect(rc RoutingContext) error {
	return rc.JsonWrap(serverConfig, nil)
}

func pathOf(path string) string {
	return fmt.Sprintf("%s%s", serverConfig.ContextPath, path)
}
func (e *EchoNi) customLog() {
	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("${time_rfc3339} | ${level} | ${short_file} |")
	}
	if serverConfig.EnableLogger {
		e.Use(middleware.Logger())
	}
	switch serverConfig.LogLevel {
	case "debug":
		e.Logger.SetLevel(log.DEBUG)
	case "info":
		e.Logger.SetLevel(log.INFO)
	case "warn":
		e.Logger.SetLevel(log.WARN)
	case "off":
		e.Logger.SetLevel(log.OFF)
	default:
		e.Logger.SetLevel(log.ERROR)
	}
}

// LoadServerCofing...
func loadServerConfig() *ServerConfig {
	return gorux.LoadConfigFile("config/server.json", &ServerConfig{
		Name:         "Biz WebAPI Server",
		Developer:    "Steven Chen",
		Version:      "v1.0.0",
		Branch:       "dev",
		LogLevel:     "off",
		EnableLogger: false,
	}).(*ServerConfig)
}
