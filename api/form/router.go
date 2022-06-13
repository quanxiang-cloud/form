package api

import (
	"github.com/quanxiang-cloud/form/internal/service/form"
	"github.com/quanxiang-cloud/form/pkg/misc/client"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"
	"github.com/quanxiang-cloud/form/pkg/misc/probe"

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
	internalHome = "internalHome"
)

// Router routing.
type Router struct {
	c *config2.Config

	engine *gin.Engine

	engineInner *gin.Engine

	Probe *probe.Probe
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
	appCenterClient := client.NewAppCenterClient(c)
	r := map[string]*gin.RouterGroup{
		managerPath:  engine.Group("/api/v1/form/:appID/m", appCenterClient.CheckIsAppAdmin),
		homePath:     engine.Group("/api/v1/form/:appID/home"),
		v2HomePath:   engine.Group("/api/v2/form/:appID/home"),
		internalPath: engineInner.Group("/api/v1/form/:appID/internal"),
		internalHome: engineInner.Group("/api/v1/form/:appID/home"),
	}
	for _, f := range routers {
		err = f(c, r)
		if err != nil {
			return nil, err
		}
	}
	probe := probe.New()
	{
		engine.GET("liveness", func(c *gin.Context) {
			probe.LivenessProbe(c.Writer, c.Request)
		})

		engine.Any("readiness", func(c *gin.Context) {
			probe.ReadinessProbe(c.Writer, c.Request)
		})
	}
	return &Router{
		c:           c,
		engine:      engine,
		engineInner: engineInner,
		Probe:       probe,
	}, nil
}

func permitRouter(c *config2.Config, r map[string]*gin.RouterGroup) error {
	permits, err := NewPermit(c)
	if err != nil {
		return err
	}
	role := r[managerPath].Group("/apiRole")
	{
		role.POST("/create", permits.CreateRole)
		role.POST("/update", permits.UpdateRole)
		role.POST("/get/:id", permits.GetRole)
		role.POST("/delete/:id", permits.DeleteRole)
		role.POST("/find", permits.FindRole)
		role.POST("/grant/list/:roleID", permits.FindGrantRole)
		role.POST("/grant/assign/:roleID", permits.AssignRoleGrant)
		role.POST("/copy", permits.CopyRole)
	}
	apiPermit := r[managerPath].Group("/apiPermit")
	{
		apiPermit.POST("/create", permits.CratePermit)
		apiPermit.POST("/update/:id", permits.UpdatePermit)
		apiPermit.POST("/get", permits.GetPermit)
		apiPermit.POST("/delete", permits.DeletePermit)
		apiPermit.POST("/list", permits.ListPermit)
	}
	home := r[homePath].Group("/apiRole") //
	{
		home.POST("/userRole/create", permits.CreateUserRole)
		home.POST("/list", permits.ListAndSelect)
		r[homePath].POST("/apiPermit/list", permits.PathPermit)
	}

	// inner 接口，permit 调用
	{
		r[internalPath].POST("/apiRole/userRole/get", permits.GetUserRole)
		r[internalPath].POST("/apiPermit/find", permits.FindPermit)
		r[internalPath].POST("/apiPermit/get", permits.GetPermit)
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
	innerHome := r[internalHome].Group("/form/:tableName")
	guide, err := form.NewRefs(c)
	if err != nil {
		return err
	}
	{
		cometHome.POST("/:action", action(guide))

		cometHome.POST("/:action/batch", batchCreate(guide))

		inner.POST("/:action", action(guide))     // inner use。
		innerHome.POST("/:action", action(guide)) // poly use

		v2Path.GET("/:id", get(guide))
		v2Path.DELETE("/:id", delete(guide))
		v2Path.PUT("/:id", update(guide))
		v2Path.POST("", create(guide))
		v2Path.GET("", search(guide))

		cometHome.GET("", search(guide))
		cometHome.GET("/relation", relation(guide))

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
	r[internalPath].POST("/schema/:tableName", table.GetTable)
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
