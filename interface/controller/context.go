package controller

import "github.com/labstack/echo/v4"

// Context abstracts the Echo context.
// This enables controller testing without a real HTTP server.
type Context interface {
	JSON(code int, i interface{}) error
	Bind(i interface{}) error
	Param(name string) string
	QueryParam(name string) string
	Status(code int) error
	Get(key string) any
	Set(key string, val any)
}

// echoContext wraps Echo's concrete Context to satisfy the controller.Context interface.
type echoContext struct {
	c echo.Context
}

func NewEchoContext(c echo.Context) Context {
	return &echoContext{c}
}

func (ec *echoContext) JSON(code int, i interface{}) error {
	return ec.c.JSON(code, i)
}

func (ec *echoContext) Bind(i interface{}) error {
	return ec.c.Bind(i)
}

func (ec *echoContext) Param(name string) string {
	return ec.c.Param(name)
}

func (ec *echoContext) QueryParam(name string) string {
	return ec.c.QueryParam(name)
}

func (ec *echoContext) Status(code int) error {
	return ec.c.NoContent(code)
}

func (ec *echoContext) Get(key string) any {
	return ec.c.Get(key)
}

func (ec *echoContext) Set(key string, val any) {
	ec.c.Set(key, val)
}
