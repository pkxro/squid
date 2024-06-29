package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkxro/squid/internal/common"
	"github.com/pkxro/squid/internal/controller"
	"github.com/pkxro/squid/internal/model"
	"github.com/pkxro/squid/middleware"
	"github.com/pkxro/squid/pkg"
)

// Manager is a struct that holds the controller manager and a router
type Manager struct {
	Controller controller.Manager
	Router     *gin.Engine
}

const (
	versionRouterGroup = "/v1"
	txRouterGroup      = "/transaction"
)

// RegisterRouters is a router method to add all nested routers to the engine
func (m *Manager) RegisterRouters(router *gin.Engine) {
	// V1 API Router Group
	v1 := router.Group(versionRouterGroup)

	// Tx Router Group
	m.RegisterTransactionsRouter(v1.Group(txRouterGroup))

}

// NewRouterManager is a constructor that returns a new instance of RouterManager
func NewRouterManager(
	cm controller.Manager,
) *Manager {

	gin.SetMode(gin.ReleaseMode)

	if (cm.Ocfg.Environment == model.ApplicationEnvironmentDev) ||
		(cm.Ocfg.Environment == model.ApplicationEnvironmentLocal) {
		gin.SetMode(gin.DebugMode)
	}

	return &Manager{
		Controller: cm,
		Router:     gin.New(),
	}

}

// InitRouter is a method that initializes the router
func (m *Manager) InitRouter() error {
	m.Router.Use(common.ZapJSONLogger(m.Controller.Logger))

	// Recover from panics
	m.Router.Use(gin.Recovery())

	// Cors
	m.Router.Use(cors.Default())

	// Rate limit
	m.Router.Use(middleware.RateLimitMiddleware(
		middleware.NewRateLimiter(m.Controller.Ocfg.RateLimit, int(m.Controller.Ocfg.RateLimitInterval)),
	))

	// No proxies
	err := m.Router.SetTrustedProxies(nil)
	if err != nil {
		return err
	}

	// Register all routers
	m.RegisterRouters(m.Router)

	// Health Check
	m.Router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// Default route
	m.Router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"running": true,
			"version": pkg.APIVersion,
		})
	})

	return nil
}
