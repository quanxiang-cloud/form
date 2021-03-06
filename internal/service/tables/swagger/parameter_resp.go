package swagger

import (
	"github.com/go-openapi/spec"
)

func countResp() *spec.Responses {
	respSchemas := spec.SchemaProperties{
		"total": spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "count",
				Type:        []string{"number"},
				Title:       "处理的条数",
			},
		},
	}
	return response(respSchemas)
}

func countAndEntityResp(schemas spec.SchemaProperties) *spec.Responses {
	respSchemas := spec.SchemaProperties{
		"entity": spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "",
				Type:        []string{"object"},
				Title:       "entity",
				Properties:  schemas,
			},
		},
		"total": spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "",
				Type:        []string{"number"},
				Title:       "处理的条数",
			},
		},
	}
	return response(respSchemas)
}

func entitiesResp(schemas spec.SchemaProperties) *spec.Responses {
	respSchemas := spec.SchemaProperties{
		"entities": spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "",
				Type:        []string{"array"},
				Title:       "entities",
				Items: &spec.SchemaOrArray{
					Schema: &spec.Schema{
						SchemaProps: spec.SchemaProps{
							Type:       []string{"object"},
							Properties: schemas,
						},
					},
				},
				Properties: schemas,
			},
		},
		"total": spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "",
				Type:        []string{"number"},
				Title:       "总数",
			},
		},
	}
	return response(respSchemas)
}

func entityResp(schemas spec.SchemaProperties) *spec.Responses {
	respSchemas := spec.SchemaProperties{
		"entity": spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "",
				Type:        []string{"object"},
				Title:       "entity",
				Properties:  schemas,
			},
		},
	}
	return response(respSchemas)
}

func response(schemas spec.SchemaProperties) *spec.Responses {
	resp := &spec.Responses{
		ResponsesProps: spec.ResponsesProps{
			StatusCodeResponses: map[int]spec.Response{
				200: {
					ResponseProps: spec.ResponseProps{
						Description: "200 is ok",
						Schema: &spec.Schema{
							SchemaProps: spec.SchemaProps{
								Description: "desc",
								Type:        []string{"object"},
								Properties: spec.SchemaProperties{
									"code": spec.Schema{
										SchemaProps: spec.SchemaProps{
											Type:  []string{"number"},
											Title: "状态码",
										},
									},
									"data": spec.Schema{
										SchemaProps: spec.SchemaProps{
											Description: "",
											Type:        []string{"object"},
											Title:       "数据",
											Properties:  schemas,
										},
									},
									"msg": spec.Schema{
										SchemaProps: spec.SchemaProps{
											Type:  []string{"string"},
											Title: "描述信息",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return resp
}

func entityParameter(schemas spec.SchemaProperties, require []string) spec.Parameter {
	return spec.Parameter{
		ParamProps: spec.ParamProps{
			Name:        "root",
			In:          "body",
			Description: "body in inputs",
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Title:      "entity",
					Type:       []string{"object"},
					Properties: schemas,
					Required:   require,
				},
			},
		},
	}
}

func idParameter() spec.Parameter {
	return spec.Parameter{
		SimpleSchema: spec.SimpleSchema{
			Type: "string",
		}, ParamProps: spec.ParamProps{
			Name:        "id",
			In:          "path",
			Description: "id of the order that needs to be deleted",
			Required:    true,
		},
	}
}

func queryParameter() spec.Parameter {
	return spec.Parameter{
		SimpleSchema: spec.SimpleSchema{
			Type: "string",
		}, ParamProps: spec.ParamProps{
			Name:        "query",
			In:          "query",
			Description: "query ParamProps",
			Required:    false,
		},
	}
}

func v1QueryID() spec.Parameter {
	return spec.Parameter{
		ParamProps: spec.ParamProps{
			Name:        "root",
			In:          "body",
			Description: "body in inputs",
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: []string{"object"},
					Properties: map[string]spec.Schema{
						"query": {
							SchemaProps: getIDPar(),
						},
					},
					Required: []string{"query"},
				},
			},
		},
	}
}

func v1EntityAndID(schemas spec.SchemaProperties, require []string) spec.Parameter {
	return spec.Parameter{
		ParamProps: spec.ParamProps{
			Name:        "root",
			In:          "body",
			Description: "body in inputs",
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: []string{"object"},
					Properties: map[string]spec.Schema{
						"query": {
							SchemaProps: getIDPar(),
						},
						"entity": {
							SchemaProps: spec.SchemaProps{
								Type:       []string{"object"},
								Properties: schemas,
								Required:   require,
							},
						},
					},
					Required: []string{"query"},
				},
			},
		},
	}
}

func getIDPar() spec.SchemaProps {
	return spec.SchemaProps{
		Type: []string{"object"},
		Properties: spec.SchemaProperties{
			"term": {
				SchemaProps: spec.SchemaProps{
					Type: []string{"object"},
					Properties: spec.SchemaProperties{
						"_id": {
							SchemaProps: spec.SchemaProps{
								Type:       []string{"string"},
								Properties: spec.SchemaProperties{},
							},
						},
					},
					Required: []string{"_id"},
				},
			},
		},
		Required: []string{"term"},
	}
}

func v1OnlyEntity(schemas spec.SchemaProperties, require []string) spec.Parameter {
	return spec.Parameter{
		ParamProps: spec.ParamProps{
			Name:        "root",
			In:          "body",
			Description: "body in inputs",
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: []string{"object"},
					Properties: map[string]spec.Schema{
						"entity": {
							SchemaProps: spec.SchemaProps{
								Type:       []string{"object"},
								Properties: schemas,
								Required:   require,
							},
						},
					},
				},
			},
		},
	}
}

func v1Search() spec.Parameter {
	return spec.Parameter{
		ParamProps: spec.ParamProps{
			Name:        "root",
			In:          "body",
			Description: "body in inputs",
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: []string{"object"},
					Properties: map[string]spec.Schema{
						"query": {
							SchemaProps: spec.SchemaProps{
								Type: []string{"object"},
							},
						},
						"size": {
							SchemaProps: spec.SchemaProps{
								Type: []string{"number"},
							},
						},
						"page": {
							SchemaProps: spec.SchemaProps{
								Type: []string{"number"},
							},
						},
						"sort": {
							SchemaProps: getItem("string"),
						},
					},
					Required: []string{"query", "size", "page", "sort"},
				},
			},
		},
	}
}

func getItem(types string) spec.SchemaProps {
	return spec.SchemaProps{
		Type: []string{"array"},
		Items: &spec.SchemaOrArray{
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: []string{types},
				},
			},
		},
	}
}
