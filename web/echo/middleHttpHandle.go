package main

import (
	"github.com/labstack/echo"
	"net/http"
)

func middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println("middleware")
		h.ServeHTTP(w, r)
	})
}

func main() {
	e := echo.New()
	e.Use(echo.WrapMiddleware(middleware))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Echo!")
	})
	e.Start(":1323")
}
