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
	v2HomePath   = "v2Home"
)

// Router routing.
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
	tableRouter,
}

// NewRouter enable routing.
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
		v2HomePath:   engine.Group("/api/v2/form/:appID/home"),
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
	role := r[managerPath].Group("/apiRole")
	{
		role.POST("/create", permits.CreateRole)     //  创建角色
		role.POST("/update", permits.UpdateRole)     //  更新角色
		role.POST("/get/:id", permits.GetRole)       // 获取单条角色
		role.POST("/delete/:id", permits.DeleteRole) // 删除对应的角色
		role.POST("/find", permits.FindRole)         // 获取角色列表
		role.POST("/userRoleMatch", permits.UserRoleMatch)
		role.POST("/grant/list/:roleID", permits.FindGrantRole)     // 获取某个角色对应的人或者部门
		role.POST("/grant/assign/:roleID", permits.AssignRoleGrant) // 给某个角色加人 、减人
	}
	apiPermit := r[managerPath].Group("/apiPermit")
	{
		apiPermit.POST("/create", permits.CratePermit)      // 创建权限
		apiPermit.POST("/update/:id", permits.UpdatePermit) // 更新权限
		apiPermit.POST("/get", permits.GetPermit)           // 获取权限
		apiPermit.POST("/list", permits.ListPermit)         // 获取权限
		apiPermit.POST("/delete", permits.DeletePermit)     // 删除权限
		apiPermit.POST("/find", permits.FindPermit)         // 获取权限
	}
	home := r[homePath].Group("/permission")
	{
		home.POST("/perGroup/saveUserPerMatch", permits.SaveUserPerMatch) // 保存用户匹配的权限组
	}

	return nil
}

func cometRouter(c *config2.Config, r map[string]*gin.RouterGroup) error {
	cometHome := r[homePath].Group("/form/:tableName")
	v2Path := r[v2HomePath].Group("/form/:tableName")

	guide, err := form.NewRefs(c)
	if err != nil {
		return err
	}
	{
		cometHome.POST("/:action", action(guide))

		v2Path.GET("/:id", get(guide))
		v2Path.DELETE("/:id", delete(guide))
		v2Path.PUT("/:id", update(guide))
		v2Path.POST("", create(guide))
		v2Path.GET("", search(guide))
	}
	return nil
}

func tableRouter(c *config2.Config, r map[string]*gin.RouterGroup) error {
	table, err := NewTable(c)
	if err != nil {
		return err
	}
	manager := r[managerPath].Group("/table")
	{
		manager.POST("/create", table.CrateTable)
		manager.POST("/getByID", table.GetTable)
		manager.POST("/delete", table.DeleteTable)
		manager.POST("/createBlank", table.CreateBlank)
		manager.POST("/search", table.FindTable)
	}
	managerConfig := r[managerPath].Group("/config")
	{
		managerConfig.POST("/create", table.UpdateConfig)
	}
	return nil
}

func innerRouter(c *config2.Config, r map[string]*gin.RouterGroup) error {
	backup, err := NewBackup(c)
	if err != nil {
		return err
	}
	bg := r[internalPath].Group("/backup")
	{
		bg.POST("/export/table", backup.ExportTable)
		bg.POST("/export/permit", backup.ExportPermit)
		bg.POST("/export/role", backup.ExportRole)
		bg.POST("/export/tableSchema", backup.ExportTableSchema)
		bg.POST("/export/tableRelation", backup.ExportTableRelation)
	}
	{
		bg.POST("/import/table", backup.ImportTable)
		bg.POST("/import/permit", backup.ImportPermit)
		bg.POST("/import/role", backup.ImportRole)
		bg.POST("/import/tableSchema", backup.ImportTableSchema)
		bg.POST("/import/tableRelation", backup.ImportTableRelation)
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

// Run router.
func (r *Router) Run() {
	go r.engineInner.Run(r.c.PortInner)
	r.engine.Run(r.c.Port)
}

// Close router.
func (r *Router) Close() {
}

func (r *Router) router() {
}
