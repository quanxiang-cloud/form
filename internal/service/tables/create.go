package tables

import (
	"context"

	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/form/internal/service/tables/util"

	id2 "github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/models/mysql"
	"github.com/quanxiang-cloud/form/internal/service"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
	"gorm.io/gorm"
)

const (
	_description = "description"
	_title       = "title"
)

// 处理web Table de.
type webTable struct {
	next      Guidance
	db        *gorm.DB
	tableRepo models.TableRepo
}

func (w *webTable) Do(ctx context.Context, bus *Bus) (*DoResponse, error) {
	one, err := w.tableRepo.Get(w.db, bus.AppID, bus.TableID)
	if err != nil {
		return nil, err
	}
	if one.ID != "" {
		bus.Update = true
	}
	tables := &models.Table{
		Schema: bus.Schema,
	}
	if one.ID == "" {
		tables.ID = id2.StringUUID()
		tables.TableID = bus.TableID
		tables.AppID = bus.AppID
		tables.CreatedAt = time2.NowUnix()
		err = w.tableRepo.BatchCreate(w.db, tables)
		if err != nil {
			return nil, err
		}
	} else {
		err = w.tableRepo.Update(w.db, bus.AppID, bus.TableID, tables)
		if err != nil {
			return nil, err
		}
	}
	return w.next.Do(ctx, bus)
}

func NewWebTable(conf *config.Config) (Guidance, error) {
	db, err := service.CreateMysqlConn(conf)
	if err != nil {
		return nil, err
	}

	schema, err := newTableSchema(conf)
	if err != nil {
		return nil, err
	}
	return &webTable{
		db:        db,
		tableRepo: mysql.NewTableRepo(),
		next:      schema,
	}, nil
}

type tableSchema struct {
	db              *gorm.DB
	tableSchemaRepo models.TableSchemeRepo
	next            Guidance
}

func (t *tableSchema) Do(ctx context.Context, bus *Bus) (*DoResponse, error) {
	properties, err := util.GetMapToMap(bus.Schema, _properties)
	if err != nil {
		return nil, err
	}

	schemas, total, err := util.Convert1(properties)
	description := util.GetMapToString(bus.Schema, _description)
	title := util.GetMapToString(bus.Schema, _title)

	bus.ConvertSchemas = ConvertSchemas{
		Title:         title,
		Description:   description,
		FieldLen:      total,
		ConvertSchema: schemas,
	}
	//
	tables := &models.TableSchema{
		Title:       bus.Title,
		Schema:      bus.ConvertSchema,
		FieldLen:    bus.FieldLen,
		Description: bus.Description,
	}
	if bus.Source == 0 {
		bus.Source = models.FormSource
	}
	if !bus.Update { // create
		tables.ID = id2.StringUUID()
		tables.Source = bus.Source
		tables.AppID = bus.AppID
		tables.TableID = bus.TableID
		tables.CreatedAt = time2.NowUnix()
		tables.CreatorName = bus.UserName
		tables.CreatorID = bus.UserID
		err = t.tableSchemaRepo.BatchCreate(t.db, tables)
		if err != nil {
			return nil, err
		}
	} else {
		tables.EditorID = bus.UserID
		tables.EditorName = bus.UserName
		tables.UpdatedAt = time2.NowUnix()
		tables.EditorID = bus.UserID
		tables.EditorName = bus.UserName
		err = t.tableSchemaRepo.Update(t.db, bus.AppID, bus.TableID, tables)
		if err != nil {
			return nil, err
		}
	}
	return t.next.Do(ctx, bus)
}

func newTableSchema(conf *config.Config) (Guidance, error) {
	component, err := newComponent(conf)
	if err != nil {
		return nil, err
	}
	db, err := service.CreateMysqlConn(conf)
	if err != nil {
		return nil, err
	}

	return &tableSchema{
		db:              db,
		next:            component,
		tableSchemaRepo: mysql.NewTableSchema(),
	}, nil
}
