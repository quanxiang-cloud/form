package api

import (
	"fmt"
	"github.com/go-openapi/spec"
	"testing"
)

func TestName(t *testing.T) {

	_ = spec.OperationProps{
		ID:          "operationID",
		Security:    nil,
		Deprecated:  false,
		Description: "",
		Summary:     "",
		Produces:    []string{"application/json"},
		Consumes:    []string{"application/json"},
		Responses: &spec.Responses{
			ResponsesProps: spec.ResponsesProps{
				StatusCodeResponses: map[int]spec.Response{
					200: {
						ResponseProps: spec.ResponseProps{
							Description: "",
							Schema: &spec.Schema{
								SchemaProps: spec.SchemaProps{
									Description: "",
									Type:        []string{"object"},
									Properties: spec.SchemaProperties{
										"code": spec.Schema{
											SchemaProps: spec.SchemaProps{
												Type:  []string{"number"},
												Title: "code",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},

		Parameters: []spec.Parameter{
			{
				ParamProps: spec.ParamProps{
					Name:        "root",
					In:          "body",
					Description: "body in inputs",
					Schema: &spec.Schema{
						SchemaProps: spec.SchemaProps{},
					},
				},
			},
			{},
		},
	}
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
					License:     nil,
				},
				VendorExtensible: spec.VendorExtensible{
					Extensions: spec.Extensions{},
				},
			},
			Host:    "",
			Swagger: "2.0",
			Tags: []spec.Tag{{
				TagProps: spec.TagProps{
					Name: "",
				},
			}},
			Schemes:  []string{"http"},
			BasePath: "/",
			Consumes: []string{"application/json"},
			Produces: []string{"application/json"},
			Paths: &spec.Paths{
				Paths: make(map[string]spec.PathItem),
			},
			Definitions: make(map[string]spec.Schema),
		},
	}

	json, err := swagger.MarshalJSON()

	if err != nil {
		fmt.Println(json)
	}
	s := string(json)

	fmt.Println(s)

}

//func getSchema (schema map[string] interface{}) spec.Schema {
//
//
//	//specSchema := &spec.Schema{
//	//	SchemaProps:spec.SchemaProps{
//	//
//	//	},
//	//}
//	//props := spec.SchemaProps{
//	//
//	//}
//	//props.Title = ""
//	//for key ,value := range schema {
//	//
//	//}
//}
