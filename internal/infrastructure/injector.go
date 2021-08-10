package infrastructure

import (
	"context"

	"go.uber.org/zap"

	"just-to-do-it/internal/config"
	"just-to-do-it/internal/controllers"
	"just-to-do-it/internal/interfaces"
)

type IInjector interface {
	InjectTaskController() controllers.TaskController
}

var env *environment

type environment struct {
	port string

	logger    *zap.SugaredLogger
	dbHandler interfaces.DBHandler
}

func (e *environment) InjectTaskController() controllers.TaskController {
	return controllers.TaskController{
		Log: e.logger,
	}
}

func Injector(logger *zap.SugaredLogger, ctx context.Context, cfg *config.Config) (IInjector, error) {
	dbClient, err := initPostgresClient(cfg, ctx)
	if err != nil {
		return nil, err
	}

	env = &environment{
		port: cfg.Port,

		logger:    logger,
		dbHandler: dbClient,
	}

	return env, err
}
