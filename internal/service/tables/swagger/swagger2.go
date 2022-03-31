package swagger

import (
	"fmt"
	"github.com/go-openapi/spec"
)

func V1GetMethod(tableID, tableName string, schemas spec.SchemaProperties) spec.OperationProps {
	responses := entityResp(schemas)
	parameters := []spec.Parameter{
		v1QueryID(),
	}
	return doOperationProps(fmt.Sprintf("%s_get", tableID), getSummary(tableName, "查询单条v1"), responses, parameters)
}

func V1Delete(tableID, tableName string, schemas spec.SchemaProperties) spec.OperationProps {
	responses := countResp()
	parameters := []spec.Parameter{
		v1QueryID(),
	}
	return doOperationProps(fmt.Sprintf("%s_delete", tableID), getSummary(tableName, "删除v1"), responses, parameters)
}

func V1SearchMethod(tableID, tableName string, schemas spec.SchemaProperties) spec.OperationProps {
	responses := entitiesResp(schemas)
	parameters := []spec.Parameter{
		v1Search(),
	}
	return doOperationProps(fmt.Sprintf("%s_search", tableID), getSummary(tableName, "查询多条v1"), responses, parameters)
}

func V1Update(tableID, tableName string, schemas, filterSchema spec.SchemaProperties) spec.OperationProps {
	responses := countAndEntityResp(schemas)
	parameters := []spec.Parameter{
		v1EntityAndID(filterSchema),
	}
	return doOperationProps(fmt.Sprintf("%s_update", tableID), getSummary(tableName, "更新v1"), responses, parameters)
}

func V1Create(tableID, tableName string, schemas, filterSchema spec.SchemaProperties) spec.OperationProps {
	responses := countAndEntityResp(schemas)
	parameters := []spec.Parameter{
		v1OnlyEntity(filterSchema),
	}
	return doOperationProps(fmt.Sprintf("%s_create", tableID), getSummary(tableName, "创建v1"), responses, parameters)
}
