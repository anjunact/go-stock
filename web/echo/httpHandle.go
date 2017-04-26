package echo

import (
	"github.com/labstack/echo"
	"net/http"
	"io"
)

func handler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Echo!")
}

func main() {
	e := echo.New()
	e.GET("/", echo.WrapHandler(http.HandlerFunc(handler)))
	e.Start(":1323")
}
