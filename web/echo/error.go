package main

import (
"net/http"

"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract the credentials from HTTP request header and perform a security
			// check

			// For invalid credentials
			return echo.NewHTTPError(http.StatusUnauthorized)

			// For valid credentials call next
			// return next(c)
		}
	})
	e.GET("/", welcome)
	e.Logger.Fatal(e.Start(":1323"))
}

func welcome(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome!")
}
