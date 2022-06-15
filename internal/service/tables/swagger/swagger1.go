package swagger

import (
	"encoding/json"
	"fmt"

	"github.com/go-openapi/spec"
	"github.com/quanxiang-cloud/form/internal/service/tables/util"
)

const (
	url1   = "/api/v1/form/%s/home/form/%s/%s"
	url2   = "/api/v2/form/%s/home/form/%s"
	url3   = "/api/v2/form/%s/home/form/%s/:id"
	get    = "get"
	create = "create"
	update = "update"
	delete = "delete"
	search = "search"
)

func GetMethod(schemasBus *schemasBus) spec.OperationProps {
	responses := entityResp(schemasBus.schemas)
	parameters := []spec.Parameter{
		idParameter(),
	}
	return doOperationProps(&operation{
		"v2_get",
		util.GetSummary(schemasBus.tableName, "查询单条"),
		responses,

		parameters,
	})
}

func PutMethod(schemasBus *schemasBus) spec.OperationProps {
	responses := countAndEntityResp(schemasBus.schemas)
	parameters := []spec.Parameter{
		//	entityParameter(schemasBus.filterSchema),
		idParameter(),
	}
	return doOperationProps(&operation{
		"v2_update",
		util.GetSummary(schemasBus.tableName, "更新"),
		responses,

		parameters,
	})
}

func DeleteMethod(schemasBus *schemasBus) spec.OperationProps {
	responses := countResp()
	parameters := []spec.Parameter{
		// entityParameter(schemasBus.schemas),
		idParameter(),
	}
	return doOperationProps(&operation{
		"v2_delete",
		util.GetSummary(schemasBus.tableName, "删除"),
		responses,

		parameters,
	})
}

func PostMethod(schemasBus *schemasBus) spec.OperationProps {
	responses := countAndEntityResp(schemasBus.schemas)

	parameters := []spec.Parameter{
		// entityParameter(schemasBus.filterSchema),
		idParameter(),
	}
	return doOperationProps(&operation{
		"v2_create",
		util.GetSummary(schemasBus.tableName, "创建"),
		responses,

		parameters,
	})
}

func SearchMethod(schemasBus *schemasBus) spec.OperationProps {
	responses := entitiesResp(schemasBus.schemas)
	parameters := []spec.Parameter{
		queryParameter(),
	}
	return doOperationProps(&operation{
		"v2_search",
		util.GetSummary(schemasBus.tableName, "查询多条"),
		responses,
		parameters,
	})
}

type operation struct {
	operationID string
	summary     string
	response    *spec.Responses
	parameter   []spec.Parameter
}

func doOperationProps(operation *operation) spec.OperationProps {
	return spec.OperationProps{
		ID:          operation.operationID,
		Security:    nil,
		Deprecated:  false,
		Description: "",
		Summary:     operation.summary,
		Produces:    []string{"application/json"},
		Consumes:    []string{"application/json"},
		Responses:   operation.response,
		Parameters:  operation.parameter,
	}
}

func DoSchemas(appID, tableID, tableName string, schemas spec.SchemaProperties, require []string) (string, error) {
	filterSystems := make(spec.SchemaProperties)
	schemasbus := &schemasBus{
		tableName:    tableName,
		tableID:      tableID,
		schemas:      schemas,
		filterSchema: filterSystems,
	}

	util.FilterSystem(schemas, filterSystems)
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
								OperationProps: SearchMethod(schemasbus),
							},
							Post: &spec.Operation{ // create
								OperationProps: PostMethod(schemasbus),
							},
						},
					},
					fmt.Sprintf(url3, appID, tableID): {
						PathItemProps: spec.PathItemProps{
							Get: &spec.Operation{ // get
								OperationProps: GetMethod(schemasbus),
							},
							Put: &spec.Operation{ //  put
								OperationProps: PutMethod(schemasbus),
							},
							Delete: &spec.Operation{ //  delete
								OperationProps: DeleteMethod(schemasbus),
							},
						},
					},
					fmt.Sprintf(url1, appID, tableID, get): {
						PathItemProps: spec.PathItemProps{
							Post: &spec.Operation{ // get
								OperationProps: V1GetMethod(schemasbus),
							},
						},
					},
					fmt.Sprintf(url1, appID, tableID, search): {
						PathItemProps: spec.PathItemProps{
							Post: &spec.Operation{ // get
								OperationProps: V1SearchMethod(schemasbus),
							},
						},
					},
					fmt.Sprintf(url1, appID, tableID, update): {
						PathItemProps: spec.PathItemProps{
							Post: &spec.Operation{ // get
								OperationProps: V1Update(schemasbus, require),
							},
						},
					},
					fmt.Sprintf(url1, appID, tableID, delete): {
						PathItemProps: spec.PathItemProps{
							Post: &spec.Operation{ // get
								OperationProps: V1Delete(schemasbus),
							},
						},
					},
					fmt.Sprintf(url1, appID, tableID, create): {
						PathItemProps: spec.PathItemProps{
							Post: &spec.Operation{ // get
								OperationProps: V1Create(schemasbus, require),
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

func DoSchemas1(appID, tableID, tableName string, schemas spec.SchemaProperties) (string, error) {
	filterSystems := make(spec.SchemaProperties)

	schemasbus := &schemasBus{
		tableName:    tableName,
		tableID:      tableID,
		schemas:      schemas,
		filterSchema: filterSystems,
	}

	util.FilterSystem(schemas, filterSystems)
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
					//fmt.Sprintf(url2, appID, tableID): {
					//	PathItemProps: spec.PathItemProps{
					//		Get: &spec.Operation{ // search
					//
					//			OperationProps: SearchMethod(schemasbus),
					//		},
					//		Post: &spec.Operation{ // create
					//			OperationProps: PostMethod(schemasbus),
					//		},
					//	},
					//},
					//fmt.Sprintf(url3, appID, tableID): {
					//	PathItemProps: spec.PathItemProps{
					//		Get: &spec.Operation{ // get
					//			OperationProps: GetMethod(schemasbus),
					//		},
					//		Put: &spec.Operation{ //  put
					//			OperationProps: PutMethod(schemasbus),
					//		},
					//		Delete: &spec.Operation{ //  delete
					//			OperationProps: DeleteMethod(schemasbus),
					//		},
					//	},
					//},
					fmt.Sprintf(url1, appID, tableID, get): {
						PathItemProps: spec.PathItemProps{
							Post: &spec.Operation{ // get
								OperationProps: V1GetMethod(schemasbus),
							},
						},
					},
					fmt.Sprintf(url1, appID, tableID, search): {
						PathItemProps: spec.PathItemProps{
							Post: &spec.Operation{ // get
								OperationProps: V1SearchMethod(schemasbus),
							},
						},
					},
					fmt.Sprintf(url1, appID, tableID, update): {
						PathItemProps: spec.PathItemProps{
							Post: &spec.Operation{ // get
								OperationProps: V1Update(schemasbus, nil),
							},
						},
					},
					fmt.Sprintf(url1, appID, tableID, delete): {
						PathItemProps: spec.PathItemProps{
							Post: &spec.Operation{ // get
								OperationProps: V1Delete(schemasbus),
							},
						},
					},
					fmt.Sprintf(url1, appID, tableID, create): {
						PathItemProps: spec.PathItemProps{
							Post: &spec.Operation{ // get
								OperationProps: V1Create(schemasbus, nil),
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
