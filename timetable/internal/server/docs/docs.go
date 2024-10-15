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
        "/api/Appointment/{id}": {
            "delete": {
                "description": "Cancel a previously made appointment",
                "tags": [
                    "Appointment"
                ],
                "summary": "Cancel an appointment",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Appointment ID",
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
        "/api/Timetable": {
            "post": {
                "description": "Create a new entry in the timetable",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Timetable"
                ],
                "summary": "Create a new timetable entry",
                "parameters": [
                    {
                        "description": "Timetable object",
                        "name": "timetable",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.timetableInfoRequestBody"
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
        "/api/Timetable/Doctor/{id}": {
            "get": {
                "description": "Retrieve the timetable for a specific doctor",
                "tags": [
                    "Timetable"
                ],
                "summary": "Get doctor's timetable by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Doctor ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "From date (ISO8601)",
                        "name": "from",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "To date (ISO8601)",
                        "name": "to",
                        "in": "query"
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
            "delete": {
                "description": "Delete all entries in the timetable for a specific doctor",
                "tags": [
                    "Timetable"
                ],
                "summary": "Delete all timetable entries for a doctor",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Doctor ID",
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
        "/api/Timetable/Hospital/{id}": {
            "get": {
                "description": "Retrieve the timetable for a specific hospital",
                "tags": [
                    "Timetable"
                ],
                "summary": "Get hospital timetable by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Hospital ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "From date (ISO8601)",
                        "name": "from",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "To date (ISO8601)",
                        "name": "to",
                        "in": "query"
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
            "delete": {
                "description": "Delete all entries in the timetable for a specific hospital",
                "tags": [
                    "Timetable"
                ],
                "summary": "Delete all timetable entries for a hospital",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Hospital ID",
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
        "/api/Timetable/Hospital/{id}/Room/{room}": {
            "get": {
                "description": "Retrieve the timetable for a specific room in a hospital",
                "tags": [
                    "Timetable"
                ],
                "summary": "Get hospital room timetable",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Hospital ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Room Name",
                        "name": "room",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "From date (ISO8601)",
                        "name": "from",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "To date (ISO8601)",
                        "name": "to",
                        "in": "query"
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
        "/api/Timetable/{id}": {
            "put": {
                "description": "Update an existing entry in the timetable",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Timetable"
                ],
                "summary": "Update a timetable entry",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Timetable ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Timetable object",
                        "name": "timetable",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.timetableInfoRequestBody"
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
            },
            "delete": {
                "description": "Delete an existing entry in the timetable",
                "tags": [
                    "Timetable"
                ],
                "summary": "Delete a timetable entry",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Timetable ID",
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
        "/api/Timetable/{id}/Appointments": {
            "get": {
                "description": "Retrieve available appointment slots based on the timetable entry",
                "tags": [
                    "Timetable"
                ],
                "summary": "Get available appointments for a given timetable entry",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Timetable ID",
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
            "post": {
                "description": "Make an appointment for a specific slot",
                "tags": [
                    "Timetable"
                ],
                "summary": "Book an appointment",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Timetable ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Appointment request",
                        "name": "appointment",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.createAppointmentRequestBody"
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
        "server.createAppointmentRequestBody": {
            "type": "object",
            "properties": {
                "time": {
                    "type": "string"
                }
            }
        },
        "server.timetableInfoRequestBody": {
            "type": "object",
            "properties": {
                "doctorId": {
                    "type": "string"
                },
                "from": {
                    "type": "string"
                },
                "hospitalId": {
                    "type": "string"
                },
                "room": {
                    "type": "string"
                },
                "to": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8083",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Account",
	Description:      "Account API documentation",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
