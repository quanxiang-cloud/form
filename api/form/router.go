package api

import (
	"net/http"

	"github.com/quanxiang-cloud/form/internal/service/form"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"

	"github.com/gin-gonic/gin"
	gin2 "github.com/quanxiang-cloud/cabin/tailormade/gin"
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

	engineInner *gin.Engine
}

type router func(c *config2.Config, r map[string]*gin.RouterGroup) error

var routers = []router{
	cometRouter,
	innerRouter,
	permitRouter,
	tableRouter,
	dataSetRouter,
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
		role.POST("/create", permits.CreateRole)                    //  创建角色
		role.POST("/update", permits.UpdateRole)                    //  更新角色
		role.POST("/get/:id", permits.GetRole)                      // 获取单条角色
		role.POST("/delete/:id", permits.DeleteRole)                // 删除对应的角色
		role.POST("/find", permits.FindRole)                        // 获取角色列表
		role.POST("/grant/list/:roleID", permits.FindGrantRole)     // 获取某个角色对应的人或者部门
		role.POST("/grant/assign/:roleID", permits.AssignRoleGrant) // 给某个角色加人 、减人
	}
	apiPermit := r[managerPath].Group("/apiPermit")
	{
		apiPermit.POST("/create", permits.CratePermit)      // 创建权限
		apiPermit.POST("/update/:id", permits.UpdatePermit) // 更新权限
		apiPermit.POST("/get", permits.GetPermit)           // 获取权限
		apiPermit.POST("/delete", permits.DeletePermit)     // 删除权限

		apiPermit.POST("/list", permits.ListPermit) // 获取权限 前端在用
	}
	home := r[homePath].Group("/apiRole") //
	{
		home.POST("/userRole/create", permits.CreateUserRole)   // 保存用户匹配的权限组
		home.POST("/list", permits.ListAndSelect)               // 查看这个人下面，有哪些，角色
		r[homePath].POST("/apiPermit/list", permits.PathPermit) //  看这个人下有那些path 权限。
	}

	// inner 接口，permit 调用
	{
		r[internalPath].POST("/apiRole/userRole/get", permits.GetUserRole)
		r[internalPath].POST("/apiPermit/find", permits.FindPermit)
		r[internalPath].POST("/apiRole/userRole/create", permits.CreateUserRole)
		r[internalPath].POST("/apiRole/create", permits.CreateRole)
		r[internalPath].POST("/apiRole/grant/assign/:roleID", permits.AssignRoleGrant)

	}

	return nil
}

func cometRouter(c *config2.Config, r map[string]*gin.RouterGroup) error {
	cometHome := r[homePath].Group("/form/:tableName")
	v2Path := r[v2HomePath].Group("/form/:tableName")
	inner := r[internalPath].Group("/form/:tableName")
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
		inner.POST("/:action", action(guide))
	}
	table, err := NewTable(c)
	if err != nil {
		return err
	}
	// get schema
	r[homePath].POST("/schema/:tableName", table.GetTable)
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
		manager.POST("/getInfo", table.GetTableInfo)
		manager.POST("/getXName", table.GetXName)
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

	engine.Use(gin2.LoggerFunc(), gin2.RecoveryFunc())

	return engine, nil
}

func dataSetRouter(c *config2.Config, r map[string]*gin.RouterGroup) error {
	dataset, err := NewDataSet(c)
	if err != nil {
		return err
	}

	datasetHome := r[homePath].Group("", func(c *gin.Context) {
		if c.Param("appID") != "dataset" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	})
	datasetManager := r[managerPath].Group("", func(c *gin.Context) {
		if c.Param("appID") != "dataset" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	})
	{
		datasetHome.POST("/get", dataset.GetDataSet)
		datasetManager.POST("/create", dataset.CreateDataSet)
		datasetManager.POST("/get", dataset.GetDataSet)
		datasetManager.POST("/update", dataset.UpdateDataSet)
		datasetManager.POST("/getByCondition", dataset.GetByConditionSet)
		datasetManager.POST("/delete", dataset.DeleteDataSet)
	}

	return nil
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
