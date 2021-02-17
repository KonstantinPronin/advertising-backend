package middleware

import (
	"github.com/KonstantinPronin/advertising-backend/pkg/model"
	"github.com/labstack/echo"
	"net/http"
)

func ParseErrors(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		err := next(ctx)

		if err != nil {
			switch err := err.(type) {
			case *echo.HTTPError:
				return err
			case *model.InvalidArgumentError:
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			case *model.ForbiddenError:
				return echo.NewHTTPError(http.StatusForbidden, err.Error())
			case *model.NotFoundError:
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			case *model.ConflictError:
				return echo.NewHTTPError(http.StatusConflict, err.Error())
			default:
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
		}

		return nil
	}
}
