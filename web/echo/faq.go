package main

import (
	"github.com/labstack/echo"
	"net/http"
	"os"
	"net"
)
//curl --unix-socket /tmp/echo.sock http://localhost
func main()  {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	os.Remove("/tmp/echo.sock")
	l, err := net.Listen("unix", "/tmp/echo.sock")
	if err != nil {
		e.Logger.Fatal(err)
	}
	e.Listener = l
	e.Logger.Fatal(e.Start(""))

}
