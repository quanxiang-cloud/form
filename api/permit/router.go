package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	guard "github.com/quanxiang-cloud/form/internal/permit/form"
	defender "github.com/quanxiang-cloud/form/internal/permit/poly"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"
)

const (
	ployPath = "poly"
	formPath = "form"
)

type router func(c *config2.Config, r map[string]*echo.Group) error

var routers = []router{
	polyRouter,
	formRouter,
}

// Router routing
type Router struct {
	c *config2.Config

	engine *echo.Echo
}

func NewRouter(c *config2.Config) (*Router, error) {
	engine := newRouter(c)

	r := map[string]*echo.Group{
		ployPath: engine.Group("*"),
		formPath: engine.Group("/api/v1/form"),
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

	engine.Use(middleware.Logger(), middleware.Recover())

	return engine
}

func (r *Router) Run() error {
	return r.engine.Start(r.c.Port)
}

func polyRouter(c *config2.Config, r map[string]*echo.Group) error {
	cor, err := defender.NewParam(c)
	if err != nil {
		return err
	}

	group := r[ployPath]
	{
		group.Any("/api/v1/poly/*", ProxyPoly(cor))
	}
	return nil
}

func formRouter(c *config2.Config, r map[string]*echo.Group) error {
	cor, err := guard.NewAuth(c)
	if err != nil {
		return err
	}

	group := r[formPath]
	{
		group.Any("/:appID/home/form/:tableID/:action", ProxyForm(cor))
	}
	return nil
}
