package router

import (
	"github.com/labstack/echo/v4"
	"github.com/quanxiang-cloud/cabin/logger"
	guard "github.com/quanxiang-cloud/form/internal/permit/form"
	defender "github.com/quanxiang-cloud/form/internal/permit/poly"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"
	echo2 "github.com/quanxiang-cloud/form/pkg/misc/echo"
)

const (
	ployPath = "poly"
	formPath = "form"
	cache    = "cache"
)

type router func(c *config2.Config, r map[string]*echo.Group) error

var routers = []router{
	polyRouter,
	formRouter,
	perCacheRouter,
}

// Router routing.
type Router struct {
	c *config2.Config

	engine *echo.Echo
}

func NewRouter(c *config2.Config) (*Router, error) {
	engine := newRouter(c)

	r := map[string]*echo.Group{
		ployPath: engine.Group("/api/v1/polyapi"),
		formPath: engine.Group("/api/v1/form"),
		cache:    engine.Group("/cache"),
	}

	for _, f := range routers {
		if err := f(c, r); err != nil {
			return nil, err
		}
	}

	return &Router{
		c:      c,
		engine: engine,
	}, nil
}

func newRouter(c *config2.Config) *echo.Echo {
	engine := echo.New()

	engine.Use(echo2.Logger, echo2.Recover)

	return engine
}

func (r *Router) Run() error {
	return r.engine.Start(r.c.Port)
}

func polyRouter(c *config2.Config, r map[string]*echo.Group) error {
	cor, err := defender.NewParam(c)
	if err != nil {
		logger.Logger.WithName("instantiation poly cor").Error(err)
		return err
	}

	group := r[ployPath]
	{
		group.Any("/*", ProxyPoly(cor))
	}
	return nil
}

func formRouter(c *config2.Config, r map[string]*echo.Group) error {
	cor, err := guard.NewAuth(c)
	if err != nil {
		logger.Logger.WithName("instantiation form cor").Error(err)
		return err
	}
	p, err := defender.NewProxy(c, c.Endpoint.Form)
	if err != nil {
		return err
	}

	group := r[formPath]
	{
		group.Any("/*", ProxyForm(p))
		group.Any("/:appID/home/form/:tableID/:action", ProxyForm(cor))
	}
	return nil
}

// 缓存一致性 ， userID roleID.
func perCacheRouter(c *config2.Config, r map[string]*echo.Group) error {
	caches, err := NewCache(c)
	if err != nil {
		logger.Logger.WithName("instantiation cache").Error(err)
		return err
	}
	r[cache].Any("/role", caches.UserRole)
	r[cache].Any("/permit", caches.Permit)

	return nil
}
