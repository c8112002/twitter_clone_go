package handler

import (
	"github.com/c8112002/twitter_clone_go/utils"
	"github.com/labstack/echo/v4"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	e := utils.NewError(err)
	res := NewErrorResponse(e)
	c.JSON(e.HttpStatusCode, res)
}
