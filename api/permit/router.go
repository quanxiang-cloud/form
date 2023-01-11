package router

import (
	"github.com/labstack/echo/v4"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/form/internal/permit/side"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"
	echo2 "github.com/quanxiang-cloud/form/pkg/misc/echo"
	"github.com/quanxiang-cloud/form/pkg/misc/probe"
)

const (
	ployPath   = "poly"
	formPath   = "form"
	cache      = "cache"
	v2FormPath = "v2Form"
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

	Probe *probe.Probe
}

func NewRouter(c *config2.Config) (*Router, error) {
	engine := newRouter(c)

	r := map[string]*echo.Group{
		ployPath:   engine.Group("/api/v1/polyapi"),
		formPath:   engine.Group("/api/v1/form"),
		cache:      engine.Group("/cache"),
		v2FormPath: engine.Group("/api/v2/form"),
	}

	for _, f := range routers {
		if err := f(c, r); err != nil {
			return nil, err
		}
	}
	probe := probe.New()
	{
		engine.GET("liveness", func(c echo.Context) error {
			probe.LivenessProbe(c.Response(), c.Request())
			return nil
		})
		engine.Any("readiness", func(c echo.Context) error {
			probe.ReadinessProbe(c.Response(), c.Request())
			return nil
		})

	}
	return &Router{
		c:      c,
		engine: engine,
		Probe:  probe,
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
	cor, err := side.NewAuth(c, c.Endpoint.Poly)
	if err != nil {
		logger.Logger.WithName("instantiation poly cor").Error(err)
		return err
	}
	//p, err := side.NewNilModifyProxy(c, c.Endpoint.Poly)
	//if err != nil {
	//	return err
	//}

	group := r[ployPath]
	{
		group.Any("/request/system/app/:appID/*", Permit(cor))
		group.Any("/request/system/app/:appID/raw/inner/form/*", Permit(cor)) // 对于form
	}
	return nil
}

func formRouter(c *config2.Config, r map[string]*echo.Group) error {
	cor, err := side.NewAuth(c, c.Endpoint.Form)
	if err != nil {
		logger.Logger.WithName("instantiation form cor").Error(err)
		return err
	}
	p, err := side.NewNilModifyProxy(c, c.Endpoint.Form)
	if err != nil {
		return err
	}

	group := r[formPath]
	{
		group.Any("/*", Permit(p))
		group.Any("/:appID/home/form/:tableID/:action", Permit(cor))
	}
	v2Form := r[v2FormPath]
	{
		v2Form.GET("/:appID/home/form/:tableID/:id", Permit(cor), V2FormPath)
		v2Form.DELETE("/:appID/home/form/:tableID/:id", Permit(cor), V2FormPath)
		v2Form.PUT("/:appID/home/form/:tableID/:id", Permit(cor), V2FormPath)
		v2Form.POST("/:appID/home/form/:tableID", Permit(cor))
		v2Form.GET("/:appID/home/form/:tableID", Permit(cor))
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
