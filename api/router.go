package api

import (
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/form/internal/service/guidance"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"
)

const (
	// DebugMode indicates mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates mode is release.
	ReleaseMode = "release"
)

const (
	managerPath  = "manager"
	homePath     = "home"
	internalPath = "internal"
)

// Router routing
type Router struct {
	c *config2.Config

	engine *gin.Engine

	// for the interaction between the process engine and the form,
	// the port is opened separately, and the verification logic is different
	engineInner *gin.Engine
}

type router func(c *config2.Config, r map[string]*gin.RouterGroup) error

var routers = []router{
	cometRouter,
	innerRouter,
	permitRouter,
}

// Newrouter enable routing
func NewRouter(c *config2.Config) (*Router, error) {
	engine, err := newRouter(c)
	if err != nil {
		return nil, err
	}
	engineInner, err := newRouter(c)
	if err != nil {
		return nil, err
	}

	r := map[string]*gin.RouterGroup{
		managerPath:  engine.Group("/api/v1/form/:appID/m"),
		homePath:     engine.Group("/api/v1/form/:appID/home"),
		internalPath: engineInner.Group("/api/v1/form/:appID/internal"),
	}
	for _, f := range routers {
		err = f(c, r)
		if err != nil {
			return nil, err
		}
	}

	return &Router{
		c:           c,
		engine:      engine,
		engineInner: engineInner,
	}, nil
}

func permitRouter(c *config2.Config, r map[string]*gin.RouterGroup) error {
	permits, err := NewPermit(c)
	if err != nil {
		return err
	}
	manager := r[managerPath].Group("/permit")

	{
		manager.POST("/role/create", permits.CreateRole) //  创建权限组
		manager.POST("/role/update", permits.UpdateRole) //  更新权限组
		manager.POST("/role/addOwner", permits.AddToRole)
		manager.POST("/role/deleteOwner", permits.DeleteOwner)
		manager.POST("/role/delete", permits.DeleteRole)
		manager.POST("/apiPermit/create", permits.CratePermit)
		manager.POST("/apiPermit/update", permits.UpdatePermit)
		manager.POST("/apiPermit/get", permits.GetPermit)
	}
	home := r[homePath].Group("/permission")
	{
		home.POST("/perGroup/saveUserPerMatch", permits.SaveUserPerMatch) // 保存用户匹配的权限组
	}

	return nil
}

func cometRouter(c *config2.Config, r map[string]*gin.RouterGroup) error {
	cometHome := r[homePath].Group("/form/:tableName")
	{
		g, err := guidance.New(c)
		if err != nil {
			return err
		}

		cometHome.POST("/:action", action(g))
		cometHome.GET("data/:id", get(g))
		cometHome.POST("data", create(g))
		cometHome.PATCH("data/:id", update(g))
		cometHome.DELETE("data/:id", delete(g))
	}
	return nil
}

func innerRouter(c *config2.Config, r map[string]*gin.RouterGroup) error {
	return nil
}

func newRouter(c *config2.Config) (*gin.Engine, error) {
	if c.Model == "" || (c.Model != ReleaseMode && c.Model != DebugMode) {
		c.Model = ReleaseMode
	}
	gin.SetMode(c.Model)
	engine := gin.New()

	engine.Use(gin.Logger(),
		gin.Recovery())

	return engine, nil
}

// Run router
func (r *Router) Run() {
	go r.engineInner.Run(r.c.PortInner)
	r.engine.Run(r.c.Port)
}

// Close router
func (r *Router) Close() {
}

func (r *Router) router() {
}
