package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type CacheControlConfig struct {
	MaxAge         time.Duration
	Public         bool
	MustRevalidate bool
}

var DefaultCacheControlConfig = CacheControlConfig{
	MaxAge:         1 * time.Hour,
	Public:         true,
	MustRevalidate: true,
}

func CacheControl() echo.MiddlewareFunc {
	return CacheControlWithConfig(DefaultCacheControlConfig)
}

func CacheControlWithConfig(config CacheControlConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := next(c); err != nil {
				// if an error occurs, just propagate it
				// and don't set cache headers for error responses
				// or, depending on policy, you might want to set "no-cache"
				return err
			}

			// Do not cache error responses
			if c.Response().Status >= 400 {
				return nil
			}

			res := c.Response()
			headerVal := ""

			if config.Public {
				headerVal += "public, "
			} else {
				headerVal += "private, "
			}

			headerVal += "max-age=" + strconv.Itoa(int(config.MaxAge.Seconds()))

			if config.MustRevalidate {
				headerVal += ", must-revalidate"
			}

			res.Header().Set(echo.HeaderCacheControl, headerVal)
			return nil
		}
	}
}
