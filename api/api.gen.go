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

// Metric defines model for Metric.
type Metric struct {
	Cpu    *string `json:"cpu,omitempty"`
	Mem    *int    `json:"mem,omitempty"`
	Ram    *string `json:"ram,omitempty"`
	Uptime *int    `json:"uptime,omitempty"`
}

// Server defines model for Server.
type Server struct {
	// Cpu core
	Cpu *int `json:"cpu,omitempty"`

	// Disk MB
	Disk *int                `json:"disk,omitempty"`
	Id   *openapi_types.UUID `json:"id,omitempty"`

	// Memory MB
	Memory *int    `json:"memory,omitempty"`
	Title  *string `json:"title,omitempty"`
}

// SignInJSONBody defines parameters for SignIn.
type SignInJSONBody struct {
	Password *string `json:"password,omitempty"`
	Username *string `json:"username,omitempty"`
}

// SignUpJSONBody defines parameters for SignUp.
type SignUpJSONBody struct {
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
	Username *string `json:"username,omitempty"`
}

// SignInJSONRequestBody defines body for SignIn for application/json ContentType.
type SignInJSONRequestBody SignInJSONBody

// SignUpJSONRequestBody defines body for SignUp for application/json ContentType.
type SignUpJSONRequestBody SignUpJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// User login
	// (POST /auth/signin)
	SignIn(ctx echo.Context) error
	// User logout
	// (GET /auth/signout)
	SignOut(ctx echo.Context) error
	// User registration
	// (POST /auth/signup)
	SignUp(ctx echo.Context) error

	// (GET /metrics/{server_id})
	GetServerMetrics(ctx echo.Context, serverId openapi_types.UUID) error

	// (GET /reboot/{server_id})
	RebootServer(ctx echo.Context, serverId openapi_types.UUID) error

	// (DELETE /rent/{server_id})
	UnRentServer(ctx echo.Context, serverId openapi_types.UUID) error

	// (POST /rent/{server_id})
	RentServer(ctx echo.Context, serverId openapi_types.UUID) error

	// (GET /servers/available)
	GetAvailableServers(ctx echo.Context) error

	// (GET /servers/my)
	GetMyServers(ctx echo.Context) error

	// (GET /servers/{server_id})
	GetServer(ctx echo.Context, serverId openapi_types.UUID) error

	// (GET /vmusers/add/{server_id})
	CreateUserOnVm(ctx echo.Context, serverId openapi_types.UUID) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// SignIn converts echo context to params.
func (w *ServerInterfaceWrapper) SignIn(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.SignIn(ctx)
	return err
}

// SignOut converts echo context to params.
func (w *ServerInterfaceWrapper) SignOut(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.SignOut(ctx)
	return err
}

// SignUp converts echo context to params.
func (w *ServerInterfaceWrapper) SignUp(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.SignUp(ctx)
	return err
}

// GetServerMetrics converts echo context to params.
func (w *ServerInterfaceWrapper) GetServerMetrics(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "server_id" -------------
	var serverId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "server_id", ctx.Param("server_id"), &serverId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter server_id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetServerMetrics(ctx, serverId)
	return err
}

// RebootServer converts echo context to params.
func (w *ServerInterfaceWrapper) RebootServer(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "server_id" -------------
	var serverId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "server_id", ctx.Param("server_id"), &serverId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter server_id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.RebootServer(ctx, serverId)
	return err
}

// UnRentServer converts echo context to params.
func (w *ServerInterfaceWrapper) UnRentServer(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "server_id" -------------
	var serverId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "server_id", ctx.Param("server_id"), &serverId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter server_id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.UnRentServer(ctx, serverId)
	return err
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

// GetMyServers converts echo context to params.
func (w *ServerInterfaceWrapper) GetMyServers(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetMyServers(ctx)
	return err
}

// GetServer converts echo context to params.
func (w *ServerInterfaceWrapper) GetServer(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "server_id" -------------
	var serverId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "server_id", ctx.Param("server_id"), &serverId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter server_id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetServer(ctx, serverId)
	return err
}

// CreateUserOnVm converts echo context to params.
func (w *ServerInterfaceWrapper) CreateUserOnVm(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "server_id" -------------
	var serverId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "server_id", ctx.Param("server_id"), &serverId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter server_id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateUserOnVm(ctx, serverId)
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

	router.POST(baseURL+"/auth/signin", wrapper.SignIn)
	router.GET(baseURL+"/auth/signout", wrapper.SignOut)
	router.POST(baseURL+"/auth/signup", wrapper.SignUp)
	router.GET(baseURL+"/metrics/:server_id", wrapper.GetServerMetrics)
	router.GET(baseURL+"/reboot/:server_id", wrapper.RebootServer)
	router.DELETE(baseURL+"/rent/:server_id", wrapper.UnRentServer)
	router.POST(baseURL+"/rent/:server_id", wrapper.RentServer)
	router.GET(baseURL+"/servers/available", wrapper.GetAvailableServers)
	router.GET(baseURL+"/servers/my", wrapper.GetMyServers)
	router.GET(baseURL+"/servers/:server_id", wrapper.GetServer)
	router.GET(baseURL+"/vmusers/add/:server_id", wrapper.CreateUserOnVm)

}
