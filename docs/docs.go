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
        "/auth/login": {
            "post": {
                "description": "Авторизация пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Login user",
                "parameters": [
                    {
                        "description": "User credentials",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_go-park-mail-ru_2025_1_ChillGuys_internal_models.UserLoginRequestDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешный ответ",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Ошибка валидации",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "401": {
                        "description": "Неверные email или пароль",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "post": {
                "description": "Выход пользователя",
                "tags": [
                    "auth"
                ],
                "summary": "Logout user",
                "responses": {
                    "200": {
                        "description": "Успешный выход",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Создает нового пользователя, хеширует пароль и устанавливает JWT-токен в куки",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register user",
                "parameters": [
                    {
                        "description": "Данные для регистрации",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_go-park-mail-ru_2025_1_ChillGuys_internal_models.UserRegisterRequestDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Регистрация успешна (тело пустое, токен в куки)",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос",
                        "schema": {
                            "$ref": "#/definitions/github_com_go-park-mail-ru_2025_1_ChillGuys_internal_transport_utils.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "Пользователь уже существует",
                        "schema": {
                            "$ref": "#/definitions/github_com_go-park-mail-ru_2025_1_ChillGuys_internal_transport_utils.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/github_com_go-park-mail-ru_2025_1_ChillGuys_internal_transport_utils.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/me": {
            "get": {
                "description": "Получение информации о текущем пользователе",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get user info",
                "responses": {
                    "200": {
                        "description": "Информация о пользователе",
                        "schema": {
                            "$ref": "#/definitions/github_com_go-park-mail-ru_2025_1_ChillGuys_internal_models.User"
                        }
                    },
                    "401": {
                        "description": "Неверный токен",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_go-park-mail-ru_2025_1_ChillGuys_internal_models.User": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "phoneNumber": {
                    "type": "string"
                },
                "surname": {
                    "type": "string"
                }
            }
        },
        "github_com_go-park-mail-ru_2025_1_ChillGuys_internal_models.UserLoginRequestDTO": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "github_com_go-park-mail-ru_2025_1_ChillGuys_internal_models.UserRegisterRequestDTO": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "surname": {
                    "type": "string"
                }
            }
        },
        "github_com_go-park-mail-ru_2025_1_ChillGuys_internal_transport_utils.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "ChillGuys API",
	Description:      "API for ChillGuys marketplace",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
