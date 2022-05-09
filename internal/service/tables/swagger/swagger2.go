package swagger

import (
	"fmt"

	"github.com/go-openapi/spec"
	"github.com/quanxiang-cloud/form/internal/service/tables/util"
)

type schemasBus struct {
	tableID      string
	tableName    string
	schemas      spec.SchemaProperties
	filterSchema spec.SchemaProperties
}

func V1GetMethod(schemasBus *schemasBus) spec.OperationProps {
	responses := entityResp(schemasBus.schemas)
	parameters := []spec.Parameter{
		v1QueryID(),
	}
	return doOperationProps(&operation{
		fmt.Sprintf("%s_%s", schemasBus.tableID, get),
		util.GetSummary(schemasBus.tableName, "查询单条v1"),
		responses,

		parameters,
	})
}

func V1Delete(schemasBus *schemasBus) spec.OperationProps {
	responses := countResp()
	parameters := []spec.Parameter{
		v1QueryID(),
	}
	return doOperationProps(&operation{
		fmt.Sprintf("%s_%s", schemasBus.tableID, delete),
		util.GetSummary(schemasBus.tableName, "删除v1"),
		responses,

		parameters,
	})
}

func V1SearchMethod(schemasBus *schemasBus) spec.OperationProps {
	responses := entitiesResp(schemasBus.schemas)
	parameters := []spec.Parameter{
		v1Search(),
	}
	return doOperationProps(&operation{
		fmt.Sprintf("%s_%s", schemasBus.tableID, search),
		util.GetSummary(schemasBus.tableName, "查询多条v1"),
		responses,

		parameters,
	})
}

func V1Update(schemasBus *schemasBus, require []string) spec.OperationProps {
	responses := countAndEntityResp(schemasBus.schemas)
	parameters := []spec.Parameter{
		v1EntityAndID(schemasBus.filterSchema, require),
	}
	return doOperationProps(&operation{
		fmt.Sprintf("%s_%s", schemasBus.tableID, update),
		util.GetSummary(schemasBus.tableName, "更新v1"),
		responses,

		parameters,
	})
}

func V1Create(schemasBus *schemasBus, require []string) spec.OperationProps {
	responses := countAndEntityResp(schemasBus.schemas)
	parameters := []spec.Parameter{
		v1OnlyEntity(schemasBus.filterSchema, require),
	}
	return doOperationProps(&operation{
		fmt.Sprintf("%s_%s", schemasBus.tableID, create),
		util.GetSummary(schemasBus.tableName, "创建v1"),
		responses,

		parameters,
	})
}
