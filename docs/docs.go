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
        "/cashier/{cashierId}/login": {
            "post": {
                "description": "login Cashier",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authentication"
                ],
                "summary": "login cashier",
                "parameters": [
                    {
                        "type": "string",
                        "description": "cashier Id",
                        "name": "cashierId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.Passcode"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Response.WebResponse-dto_Token"
                        }
                    }
                }
            }
        },
        "/cashier/{cashierId}/logout": {
            "post": {
                "description": "logout Cashier",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authentication"
                ],
                "summary": "logout cashier",
                "parameters": [
                    {
                        "type": "string",
                        "description": "cashier Id",
                        "name": "cashierId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.Passcode"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Response.WebResponse-string"
                        }
                    }
                }
            }
        },
        "/cashier/{cashierId}/passcode": {
            "post": {
                "description": "passcode Cashier",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authentication"
                ],
                "summary": "passcode cashier",
                "parameters": [
                    {
                        "type": "string",
                        "description": "cashier Id",
                        "name": "cashierId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.Passcode"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Response.WebResponse-dto_Passcode"
                        }
                    }
                }
            }
        },
        "/cashiers": {
            "get": {
                "description": "get cashier list",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cashiers"
                ],
                "summary": "get cashier list",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Response.WebResponse-array_Model_Cashier"
                        }
                    }
                }
            },
            "post": {
                "description": "create cashier",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cashiers"
                ],
                "summary": "create cashier",
                "parameters": [
                    {
                        "description": "request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Model.Cashier"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/Response.WebResponse-string"
                        }
                    }
                }
            }
        },
        "/cashiers/{cashierId}": {
            "get": {
                "description": "get cashier detail",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cashiers"
                ],
                "summary": "get cashier detail",
                "parameters": [
                    {
                        "type": "string",
                        "description": "cashier id",
                        "name": "cashierId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Response.WebResponse-dto_CashierDetails"
                        }
                    }
                }
            },
            "put": {
                "description": "update cashier",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cashiers"
                ],
                "summary": "get update cashier",
                "parameters": [
                    {
                        "description": "request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Model.Cashier"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Response.WebResponse-string"
                        }
                    }
                }
            },
            "delete": {
                "description": "delete cashier",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cashiers"
                ],
                "summary": "delete cashier",
                "parameters": [
                    {
                        "type": "string",
                        "description": "cashier id",
                        "name": "cashierId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Response.WebResponse-string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "Model.Cashier": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "passcode": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "Response.WebResponse-array_Model_Cashier": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/Model.Cashier"
                    }
                },
                "message": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "Response.WebResponse-dto_CashierDetails": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/dto.CashierDetails"
                },
                "message": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "Response.WebResponse-dto_Passcode": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/dto.Passcode"
                },
                "message": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "Response.WebResponse-dto_Token": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/dto.Token"
                },
                "message": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "Response.WebResponse-string": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "dto.CashierDetails": {
            "type": "object",
            "properties": {
                "cashier_id": {
                    "type": "integer"
                },
                "created_at": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "dto.Passcode": {
            "type": "object",
            "properties": {
                "passcode": {
                    "type": "string"
                }
            }
        },
        "dto.Token": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8000",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Sales api docs",
	Description:      "This is the complete api documentation for sales api",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
