package controllers

import (
	http "just-to-do-it/internal/http_utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type TaskController struct {
	Log *zap.SugaredLogger
}

func (c *TaskController) HelloWorld(ctx echo.Context) error {
	return http.SuccessResponse(ctx, "Hello, World")
}
