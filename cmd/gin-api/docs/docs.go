// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "STones_",
            "url": "http://www.swagger.io/support",
            "email": "yelei@3k.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/user": {
            "get": {
                "tags": [
                    "UserQueryAPI"
                ],
                "summary": "根据用户名或用户状态查询用户信息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "用户名",
                        "name": "user_name",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "用户状态(1，正常; 2，失效)",
                        "name": "status",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "模糊查询",
                        "name": "query_value",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "用户列表",
                        "schema": {
                            "$ref": "#/definitions/schema.ListResult"
                        }
                    },
                    "400": {
                        "description": "错误信息",
                        "schema": {
                            "$ref": "#/definitions/schema.ErrorResult"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "schema.ErrorItem": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "schema.ErrorResult": {
            "type": "object",
            "properties": {
                "error": {
                    "$ref": "#/definitions/schema.ErrorItem"
                }
            }
        },
        "schema.ListResult": {
            "type": "object",
            "properties": {
                "list": {}
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.1",
	Host:             "127.0.0.1:8081",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "k3gin",
	Description:      "RBAC scaffolding based on GIN + GORM + WIRE.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
