package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/olekukonko/tablewriter"
	id2 "github.com/quanxiang-cloud/cabin/id"
	"github.com/quanxiang-cloud/cabin/logger"
	mongo2 "github.com/quanxiang-cloud/cabin/tailormade/db/mongo"
	mysql2 "github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

var (
	configPath = flag.String("config", "configs/config.yml", "-config 配置文件地址")
)

//处理旧版本的已经保存的过滤规则
func main() {
	conf, err := config.NewConfig(*configPath)
	if err != nil {
		panic(err)
	}
	client, err := mongo2.New(&conf.Mongo)
	if err != nil {
		panic(err)
	}
	db := client.Database("structor")

	m := bson.M{}
	opts := &options.FindOptions{}
	ctx := context.Background()
	errTableMsg := make([][]string, 0)
	// table
	tables := make([]Table, 0)
	cursor, err := db.Collection("table_schema").Find(ctx, m, opts)
	err = cursor.All(ctx, &tables)
	if err != nil {
		panic(err)
	}

	mysqlDB, err := mysql2.New(conf.Mysql, logger.Logger)
	if err != nil {
		panic(err)
	}
	//==============================
	var (
		totalTable   = 0
		successTable = 0
		failTable    = 0
	)
	totalTable = len(tables)
	for _, value := range tables {
		table := &models.Table{
			ID:        id2.StringUUID(),
			TableID:   value.TableID,
			AppID:     value.AppID,
			Schema:    value.Schema,
			Config:    value.Config,
			CreatedAt: time2.NowUnix(),
		}

		err = mysqlDB.Table("table").Create(table).Error
		if err != nil {
			failTable++
			msg := []string{"table_schema", value.ID, err.Error()}
			errTableMsg = append(errTableMsg, msg)
			continue
		}
		successTable++
	}
	//==============================
	var (
		totalSubTable   = 0
		successSubTable = 0
		failSubTable    = 0
	)
	subTable := make([]SubTable, 0)
	subTableCur, err := db.Collection("sub_table_relation").Find(ctx, m, opts)
	err = subTableCur.All(ctx, &subTable)
	if err != nil {
		panic(err)
	}
	totalSubTable = len(subTable)
	for _, value := range subTable {
		tableRelation := &models.TableRelation{
			ID:         id2.StringUUID(),
			TableID:    value.TableID,
			AppID:      value.AppID,
			FieldName:  value.FieldName,
			SubTableID: value.SubTableID,
			Filter:     value.Filter,
		}
		if value.SubTableType == "AssociatedRecords" {
			tableRelation.SubTableType = "associated_records"
		}
		err = mysqlDB.Table("table_relation").Create(tableRelation).Error
		if err != nil {
			failSubTable++
			msg := []string{"sub_table_relation", value.ID, err.Error()}
			errTableMsg = append(errTableMsg, msg)
			continue
		}
		successSubTable++

	}

	var (
		totalTableSchema   = 0
		successTableSchema = 0
		failTableSchema    = 0
	)
	dataBase := make([]DataBaseSchema, 0)
	dataBaseCur, err := db.Collection("database_schema").Find(ctx, m, opts)
	err = dataBaseCur.All(ctx, &dataBase)
	if err != nil {
		panic(err)
	}
	totalTableSchema = len(dataBase)

	//==============================
	for _, value := range dataBase {
		tableSchema := &TableSchema{
			ID:          id2.StringUUID(),
			TableID:     value.TableID,
			AppID:       value.AppID,
			FieldLen:    value.FieldLen,
			Title:       value.Title,
			Description: value.Description,
			Source:      value.Source,
			CreatedAt:   value.CreatedAt,
			UpdatedAt:   value.UpdatedAt,
			CreatorID:   value.CreatorID,
			CreatorName: value.CreatorName,
			EditorID:    value.EditorID,
			EditorName:  value.EditorName,
		}
		err = mysqlDB.Table("table_schema").Create(tableSchema).Error
		if err != nil {
			failTableSchema++
			msg := []string{"table_schema", value.ID, err.Error()}
			errTableMsg = append(errTableMsg, msg)
			continue
		}
		successTableSchema++
	}

	logger.Logger.Infof("total: %d  success :%d ,", totalTableSchema)

	Write([]string{"table_name", "total", "success", "fail"}, [][]string{
		{"table", fmt.Sprintf("%d", totalTable), fmt.Sprintf("%d", successTable), fmt.Sprintf("%d", failTable)},
		{"table_relation", fmt.Sprintf("%d", totalSubTable), fmt.Sprintf("%d", successSubTable), fmt.Sprintf("%d", failSubTable)},
		{"table_schema", fmt.Sprintf("%d", totalTableSchema), fmt.Sprintf("%d", successTableSchema), fmt.Sprintf("%d", failTableSchema)}})
	Write([]string{"table_name", "id", "errMsg", "err"}, errTableMsg)

}

// Write Write
func Write(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}
