package main

import (
	//"net/http"

	"github.com/labstack/echo"
	"net/http"
	"io"
	"html/template"
	"github.com/anjunact/go-stock/models"
	"github.com/labstack/gommon/log"
	"path/filepath"
	"fmt"
	"strconv"
)
type CustomContext struct {
	echo.Context
}

func (c *CustomContext) Foo() {
	println("foo")
}

func (c *CustomContext) Bar() {
	println("bar")
}
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	rootPath,_:=  filepath.Abs(".")
	fmt.Println("==="+rootPath)

	e := echo.New()
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{c}
			return h(cc)
		}
	})
	e.GET("/", func(c echo.Context) error {
		cc := c.(*CustomContext)
		cc.Foo()
		cc.Bar()
		return cc.String(200, "OK")
	})

	t := &Template{
		templates: template.Must(template.ParseGlob(rootPath+"/web/public/views/*.html")),
	}
	e.Renderer = t
	e.GET("/stocks/:page", Stocks)

	e.Static("/static", rootPath+"/web/assets")

	e.Logger.Fatal(e.Start(":1323"))
}
func Stocks(c echo.Context) error {
	page := c.Param("page")
	pageInt,_  := strconv.Atoi(page)
	stocks ,err := models.Stocks(pageInt,10)
	if err !=nil{
		log.Print(err)
	}
	return c.Render(http.StatusOK, "stocks", stocks)
}