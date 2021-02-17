package internal

import (
	"github.com/KonstantinPronin/advertising-backend/internal/advertisement/delivery"
	"github.com/KonstantinPronin/advertising-backend/internal/advertisement/repository"
	"github.com/KonstantinPronin/advertising-backend/internal/advertisement/usecase"
	"github.com/KonstantinPronin/advertising-backend/pkg/infrastructure"
	"github.com/labstack/echo"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
)

type Application struct {
	e    *echo.Echo
	port string
}

func (a *Application) Start() error {
	return a.e.Start(a.port)
}

func NewApplication(
	e *echo.Echo,
	db *infrastructure.Database,
	logger *zap.Logger,
	port string) *Application {

	sanitizer := bluemonday.UGCPolicy()
	rep := repository.NewMongoDbClient(db, logger)
	uc := usecase.NewAdvertising(rep, logger)

	delivery.NewAdvertisementHandler(e, uc, sanitizer, logger)

	return &Application{
		e:    e,
		port: port,
	}
}
