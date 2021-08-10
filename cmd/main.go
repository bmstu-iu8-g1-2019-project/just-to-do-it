package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"just-to-do-it/internal/config"
	"just-to-do-it/internal/infrastructure"
)

var (
	cfg *config.Config
	ctx context.Context
	log *zap.SugaredLogger
)

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("error loading logger: %s", err)
		os.Exit(1)
		return
	}

	log = logger.Sugar()

	cfg, err = config.New()
	if err != nil {
		log.Fatalf("config init error: %s", err)
	}
	log.Infof("Config loaded:\n%+v", cfg)

	ctx = context.Background()
}

func main() {
	defer log.Sync()

	router := echo.New()
	router.Debug = true
	router.Validator = &CustomValidator{Validator: *validator.New()}

	injector, err := infrastructure.Injector(log, ctx, cfg)
	if err != nil {
		panic(err)
	}

	registerRoutes(router, injector)

	log.Info("Starting listen and serve")
	log.Fatal(router.Start(fmt.Sprintf(":%s", cfg.Port)))
}

func registerRoutes(router *echo.Echo, injector infrastructure.IInjector) {
	taskController := injector.InjectTaskController()
	{
		router.GET("/hello", taskController.HelloWorld)
	}
}

type CustomValidator struct {
	Validator validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}
