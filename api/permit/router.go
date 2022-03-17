package router

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	_, err := NewPoly(c)
	if err != nil {
		return err
	}

	group := r[ployPath]
	group.Any("*", func(c echo.Context) error {
		fmt.Println("poly")
		return nil
	})

	return nil
}

func formRouter(c *config2.Config, r map[string]*echo.Group) error {
	form, err := NewForm(c)
	if err != nil {
		return err
	}

	group := r[formPath]
	group.Any("/:appID/home/form/:tableID/:action", form.Forward)
	return nil
}
