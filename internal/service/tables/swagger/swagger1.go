package swagger

import (
	"encoding/json"
	"fmt"
	"github.com/go-openapi/spec"
)

const (
	url1 = "/api/v1/form/%s/home/form/%s/%s"
	url2 = "/api/v1/form/%s/home/form/%s"
	url3 = "/api/v1/form/%s/home/form/%s/:id"
)

func GetMethod(tableName string, schemas spec.SchemaProperties) spec.OperationProps {
	responses := entityResp(schemas)
	parameters := []spec.Parameter{
		idParameter(),
	}
	return doOperationProps("v2_get", getSummary(tableName, "查询单条"), responses, parameters)
}

func PutMethod(tableName string, schemas, filterSchema spec.SchemaProperties) spec.OperationProps {
	responses := countAndEntityResp(schemas)
	parameters := []spec.Parameter{
		entityParameter(filterSchema),
		idParameter(),
	}
	return doOperationProps("v2_update", getSummary(tableName, "更新"), responses, parameters)
}
func DeleteMethod(tableName string, schemas spec.SchemaProperties) spec.OperationProps {
	responses := countResp()
	parameters := []spec.Parameter{
		entityParameter(schemas),
		idParameter(),
	}
	return doOperationProps("v2_delete", getSummary(tableName, "删除"), responses, parameters)
}

func PostMethod(tableName string, schemas, filterSchema spec.SchemaProperties) spec.OperationProps {
	responses := countAndEntityResp(schemas)

	parameters := []spec.Parameter{
		entityParameter(filterSchema),
		idParameter(),
	}
	return doOperationProps("v2_create", getSummary(tableName, "更新"), responses, parameters)
}

func SearchMethod(tableName string, schemas spec.SchemaProperties) spec.OperationProps {
	responses := entitiesResp(schemas)
	parameters := []spec.Parameter{
		queryParameter(),
	}
	return doOperationProps("v2_search", getSummary(tableName, "查询多条"), responses, parameters)
}
func doOperationProps(operationID, summary string, response *spec.Responses, parameter []spec.Parameter) spec.OperationProps {
	return spec.OperationProps{
		ID:          operationID,
		Security:    nil,
		Deprecated:  false,
		Description: "",
		Summary:     summary,
		Produces:    []string{"application/json"},
		Consumes:    []string{"application/json"},
		Responses:   response,
		Parameters:  parameter,
	}
}

func DoSchemas(appID, tableID, tableName string, schemas spec.SchemaProperties) (string, error) {

	filterSystems := make(spec.SchemaProperties)

	filterSystem(schemas, filterSystems)
	swagger := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Info: &spec.Info{
				InfoProps: spec.InfoProps{
					Contact: &spec.ContactInfo{
						ContactInfoProps: spec.ContactInfoProps{
							Name:  "",
							URL:   "",
							Email: "",
						},
					},
					Title:       "structor",
					Version:     "last",
					Description: "表单引擎",
				},
				VendorExtensible: spec.VendorExtensible{
					Extensions: spec.Extensions{},
				},
			},
			Host:    "",
			Swagger: "2.0",
			Tags: []spec.Tag{{
				TagProps: spec.TagProps{
					Name: "table",
				},
			}},
			Schemes:  []string{"http"},
			BasePath: "/",
			Consumes: []string{"application/json"},
			Produces: []string{"application/json"},
			Paths: &spec.Paths{
				Paths: map[string]spec.PathItem{
					fmt.Sprintf(url2, appID, tableID): {
						PathItemProps: spec.PathItemProps{
							Get: &spec.Operation{ // search

								OperationProps: SearchMethod(tableName, schemas),
							},
							Post: &spec.Operation{ // create
								OperationProps: PostMethod(tableName, schemas, filterSystems),
							},
						},
					},
					fmt.Sprintf(url3, appID, tableID): {
						PathItemProps: spec.PathItemProps{
							Get: &spec.Operation{ // get
								OperationProps: GetMethod(tableName, schemas),
							},
							Put: &spec.Operation{ //  put
								OperationProps: PutMethod(tableName, schemas, filterSystems),
							},
							Delete: &spec.Operation{ //  delete
								OperationProps: DeleteMethod(tableName, schemas),
							},
						},
					},
					fmt.Sprintf(url1, appID, tableID, "get"): {
						PathItemProps: spec.PathItemProps{
							Post: &spec.Operation{ // get
								OperationProps: V1GetMethod(tableID, tableName, schemas),
							},
						},
					},
					fmt.Sprintf(url1, appID, tableID, "search"): {
						PathItemProps: spec.PathItemProps{
							Post: &spec.Operation{ // get
								OperationProps: V1SearchMethod(tableID, tableName, schemas),
							},
						},
					},
					fmt.Sprintf(url1, appID, tableID, "update"): {
						PathItemProps: spec.PathItemProps{
							Post: &spec.Operation{ // get
								OperationProps: V1Update(tableID, tableName, schemas, filterSystems),
							},
						},
					},
					fmt.Sprintf(url1, appID, tableID, "delete"): {
						PathItemProps: spec.PathItemProps{
							Post: &spec.Operation{ // get
								OperationProps: V1Delete(tableID, tableName, schemas),
							},
						},
					},
					fmt.Sprintf(url1, appID, tableID, "create"): {
						PathItemProps: spec.PathItemProps{
							Post: &spec.Operation{ // get
								OperationProps: V1Create(tableID, tableName, schemas, filterSystems),
							},
						},
					},
				},
			},
			Definitions: make(map[string]spec.Schema),
		},
	}

	marshal, err := json.Marshal(swagger)
	if err != nil {
		return "", err
	}
	swaggers := string(marshal)
	return swaggers, nil

}
