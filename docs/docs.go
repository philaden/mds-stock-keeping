// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
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
        "/api/products": {
            "get": {
                "description": "This endpoint fetches a list of all create product",
                "produces": [
                    "application/json"
                ],
                "summary": "Attempts to get all products",
                "responses": {}
            },
            "post": {
                "description": "This endpoint accept stock information as a csv file",
                "produces": [
                    "application/json"
                ],
                "summary": "update stocks",
                "responses": {}
            }
        },
        "/api/products/:sku": {
            "get": {
                "description": "This endpoint fetches a stock product by its sku",
                "produces": [
                    "application/json"
                ],
                "summary": "Attempts to get a existing product by sku",
                "responses": {}
            }
        },
        "/api/products/singleproduct": {
            "post": {
                "description": "This endpoint creates a single stock product",
                "produces": [
                    "application/json"
                ],
                "summary": "Create a new single product",
                "responses": {}
            }
        },
        "/order": {
            "post": {
                "description": "This endpoint creates an order for a product. This endpoint is meant to reduce the unit of stock for product",
                "produces": [
                    "application/json"
                ],
                "summary": "Create a new order",
                "responses": {}
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "mds-stock-keeping docs",
	Description:      "This is the api documentation for my solution.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
