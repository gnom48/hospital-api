// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/History": {
            "post": {
                "description": "Create visit record",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "History"
                ],
                "summary": "Create visit history and appointment",
                "parameters": [
                    {
                        "description": "History object",
                        "name": "history",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.createHistoryRequestBody"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Authorization header",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/api/History/Account/{id}": {
            "get": {
                "description": "Retrieve records where {pacientId} = {id}",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "History"
                ],
                "summary": "Get visit history and appointments for an account",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Patient ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Authorization header",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/api/History/{id}": {
            "get": {
                "description": "Retrieve visit and appointment details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "History"
                ],
                "summary": "Get detailed information about a visit and appointments",
                "parameters": [
                    {
                        "type": "string",
                        "description": "History ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Authorization header",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {}
            },
            "put": {
                "description": "Update visit record",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "History"
                ],
                "summary": "Update visit history and appointment",
                "parameters": [
                    {
                        "type": "string",
                        "description": "History ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated History object",
                        "name": "history",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.createHistoryRequestBody"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Authorization header",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "server.createHistoryRequestBody": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "string"
                },
                "date": {
                    "type": "string"
                },
                "doctor_id": {
                    "type": "string"
                },
                "hospital_id": {
                    "type": "string"
                },
                "patient_id": {
                    "type": "string"
                },
                "room": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "gnom48.ru/History",
	BasePath:         "/",
	Schemes:          []string{"https"},
	Title:            "History",
	Description:      "History API (Document microservice) documentation. Отвечает за историю посещений пользователя. Отправляет запросы в микросервис аккаунтов для интроспекции токена и проверки существования связанных сущностей. Отправляет запросы в микросервис больниц для проверки существования связанных сущностей.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
