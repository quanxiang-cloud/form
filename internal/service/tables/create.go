package tables

import (
	"context"
	id2 "github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/models/mysql"
	"github.com/quanxiang-cloud/form/internal/service"
	swagger2 "github.com/quanxiang-cloud/form/internal/service/tables/swagger"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
	"gorm.io/gorm"
)

// 处理web Table de
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
		tables.ID = id2.HexUUID(true)
		tables.TableID = bus.TableID
		tables.AppID = bus.AppID
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
	properties, err := getMapToMap(bus.Schema, "properties")
	if err != nil {
		return nil, err
	}

	convert, total, err := swagger2.Convert1(properties)
	description := getMapToString(bus.Schema, "description")
	title := getMapToString(bus.Schema, "title")

	bus.ConvertSchemas = ConvertSchemas{
		Title:         title,
		Description:   description,
		FieldLen:      total,
		ConvertSchema: convert,
	}
	//
	tables := &models.TableSchema{
		Title:       bus.Title,
		Schema:      bus.Schema,
		FieldLen:    bus.FieldLen,
		Description: bus.Description,
	}

	if !bus.Update { // create
		tables.ID = id2.HexUUID(true)
		tables.Source = bus.Source
		tables.AppID = bus.AppID
		tables.TableID = bus.TableID
		//table.CreatedAt = time.Now()
		tables.CreatorName = bus.UserName
		tables.CreatorID = bus.UserID
		err = t.tableSchemaRepo.BatchCreate(t.db, tables)
		if err != nil {
			return nil, err
		}
	} else {
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
