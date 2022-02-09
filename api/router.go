package api

import (
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/form/internal/service/form"
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

// Router 路由
type Router struct {
	c *config2.Config

	engine *gin.Engine

	// 流程引擎与表单的交互，单独开端口，校验逻辑不同
	engineInner *gin.Engine
}

type router func(c *config2.Config, r map[string]*gin.RouterGroup) error

var routers = []router{
	permissionRouter,
	cometRouter,
}

// NewRouter 开启路由
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

func permissionRouter(c *config2.Config, r map[string]*gin.RouterGroup) error {
	permission, err := NewPermission(c)
	if err != nil {
		return err
	}
	manager := r[managerPath].Group("/permission")
	{
		manager.POST("/perGroup/create", permission.CreatePerGroup)     //  创建权限组
		manager.POST("/perGroup/updateName", permission.UpdatePerGroup) //  更新权限组
		manager.POST("/perGroup/update", permission.AddOwnerToGroup, permission.AddOwnerToApp)

		manager.POST("/perGroup/delete", permission.DelPerGroup)   // 删除权限组 ,要删除缓存
		manager.POST("/perGroup/getByID", permission.GetPerGroup)  // 根据id 获取权限组
		manager.POST("/perGroup/getList", permission.FindPerGroup) // 根据条件获取 权限组列表

		manager.POST("/perGroup/saveForm", permission.SaveForm) // 保存表单权限

		manager.POST("/perGroup/deleteForm", permission.DeleteForm) // 删除表单权限

		manager.POST("/perGroup/getForm", permission.FindForm) // 获取form 的信息
		manager.POST("/perGroup/getPerData", permission.GetPerInfo)

		manager.POST("/perGroup/getPerGroupByMenu", permission.GetPerGroupsByMenu)
	}

	home := r[homePath].Group("/permission")
	{
		home.POST("/operatePer/getOperate", permission.GetOperate) // 跟据用户id 和 部门ID，得到操作权限
		home.POST("/perGroup/getPerOption", permission.GetPerOption)
		home.POST("/perGroup/saveUserPerMatch", permission.SaveUserPerMatch) // 保存用户匹配的权限组
	}

	return nil

}

func cometRouter(c *config2.Config, r map[string]*gin.RouterGroup) error {
	authForm, err := form.NewAuthForm(c)
	if err != nil {
		return err
	}
	// form := form.NewForm()

	cometHome := r[homePath].Group("/form/:tableName")
	{
		cometHome.POST("/search", Search(authForm))
		// cometHome.POST("/get", Get(form, true))
		//cometHome.POST("/create", Create(form ,true ))
		//cometHome.POST("/update", Update(form ,true ))
		//cometHome.POST("/delete", Delete(form ,true ))

		cometHome.POST("/:action", Action(form.NewPoly()))
	}
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
