package http_utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg,omitempty"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(ctx echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Msg:     "",
		Success: true,
		Data:    data,
	})
}

func BadRequestResponse(ctx echo.Context, msg string) error {
	return ctx.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Msg:     msg,
		Success: false,
		Data:    nil,
	})
}

func InternalServerErrorResponse(ctx echo.Context, msg string) error {
	return ctx.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Msg:     msg,
		Success: false,
		Data:    nil,
	})
}

func NotFoundResponse(ctx echo.Context, msg string) error {
	return ctx.JSON(http.StatusNotFound, Response{
		Code:    http.StatusNotFound,
		Msg:     msg,
		Success: false,
		Data:    nil,
	})
}

func UnauthorizedResponse(ctx echo.Context, msg string) error {
	return ctx.JSON(http.StatusUnauthorized, Response{
		Code:    http.StatusUnauthorized,
		Msg:     msg,
		Success: false,
		Data:    nil,
	})
}

func ForbiddenResponse(ctx echo.Context, msg string) error {
	return ctx.JSON(http.StatusForbidden, Response{
		Code:    http.StatusForbidden,
		Msg:     msg,
		Success: false,
		Data:    nil,
	})
}

func LockedResponse(ctx echo.Context, msg string) error {
	return ctx.JSON(http.StatusLocked, Response{
		Code:    http.StatusLocked,
		Msg:     msg,
		Success: false,
		Data:    nil,
	})
}
