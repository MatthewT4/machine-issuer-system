// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Server defines model for Server.
type Server struct {
	// Cpu core
	Cpu *interface{} `json:"cpu,omitempty"`

	// Disk MB
	Disk *interface{}        `json:"disk,omitempty"`
	Id   *openapi_types.UUID `json:"id,omitempty"`

	// Memory MB
	Memory *interface{} `json:"memory,omitempty"`
	Title  *string      `json:"title,omitempty"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /rent/{server_id})
	RentServer(ctx echo.Context, serverId openapi_types.UUID) error

	// (GET /servers/available)
	GetAvailableServers(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// RentServer converts echo context to params.
func (w *ServerInterfaceWrapper) RentServer(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "server_id" -------------
	var serverId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "server_id", ctx.Param("server_id"), &serverId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter server_id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.RentServer(ctx, serverId)
	return err
}

// GetAvailableServers converts echo context to params.
func (w *ServerInterfaceWrapper) GetAvailableServers(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetAvailableServers(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/rent/:server_id", wrapper.RentServer)
	router.GET(baseURL+"/servers/available", wrapper.GetAvailableServers)

}
