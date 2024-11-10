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
        "/api/Accounts": {
            "get": {
                "description": "Retrieve a list of all accounts",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Accounts"
                ],
                "summary": "Get all accounts",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Start index",
                        "name": "from",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of records",
                        "name": "count",
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
            "post": {
                "description": "Create a new user account by admin",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Accounts"
                ],
                "summary": "Create a new account",
                "parameters": [
                    {
                        "description": "Account Creation Data",
                        "name": "requestBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.createAccountRequestBody"
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
        "/api/Accounts/Me": {
            "get": {
                "description": "Retrieve the current account's data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Accounts"
                ],
                "summary": "Get current account",
                "parameters": [
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
        "/api/Accounts/Update": {
            "put": {
                "description": "Update the current account's information",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Accounts"
                ],
                "summary": "Update account",
                "parameters": [
                    {
                        "description": "Account Update Data",
                        "name": "requestBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.updateAccountRequestBody"
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
        "/api/Accounts/{id}": {
            "put": {
                "description": "Update a user account by Id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Accounts"
                ],
                "summary": "Update account by Id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Account Details",
                        "name": "requestBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.createAccountRequestBody"
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
                "description": "Soft delete a user account by ID",
                "tags": [
                    "Accounts"
                ],
                "summary": "Soft delete an account by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User Id",
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
        "/api/Authentication/Refresh": {
            "get": {
                "description": "Refresh token pair by creation token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Refresh",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization header (creation token)",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/api/Authentication/SignIn": {
            "post": {
                "description": "Authenticates a user based on their username and password and generates tokens",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Sign in a user",
                "parameters": [
                    {
                        "description": "User Credentials",
                        "name": "requestBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.signInRequestBody"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/api/Authentication/SignOut": {
            "head": {
                "description": "Delete token pair by creation token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "SignOut",
                "parameters": [
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
        "/api/Authentication/SignUp": {
            "post": {
                "description": "SignUp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "SignUp",
                "parameters": [
                    {
                        "description": "User Credentials",
                        "name": "requestBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.signUpRequestBody"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/api/Authentication/Validate": {
            "get": {
                "description": "Validate regular token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Validate",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization header",
                        "name": "AccessToken",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/api/Doctors": {
            "get": {
                "description": "Retrieve a list of doctors with optional name filtering",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Doctors"
                ],
                "summary": "Get list of doctors",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Filter by doctor's full name",
                        "name": "nameFilter",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Pagination start",
                        "name": "from",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of records per page",
                        "name": "count",
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
        "/api/Doctors/{id}": {
            "get": {
                "description": "Retrieve doctor information by doctor ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Doctors"
                ],
                "summary": "Get doctor information by ID",
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
        }
    },
    "definitions": {
        "server.createAccountRequestBody": {
            "type": "object",
            "properties": {
                "first_name": {
                    "type": "string"
                },
                "last_name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "roles": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "server.signInRequestBody": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "server.signUpRequestBody": {
            "type": "object",
            "properties": {
                "first_name": {
                    "type": "string"
                },
                "last_name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "server.updateAccountRequestBody": {
            "type": "object",
            "properties": {
                "first_name": {
                    "type": "string"
                },
                "last_name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8081",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Account",
	Description:      "Account API (Account microservice) documentation. Jтвечает за авторизацию и данные о пользователях. Все остальные сервисы зависят от него, ведьименно он выпускает JWT токен и проводит интроспекцию.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
