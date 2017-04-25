package main

import (
	"github.com/labstack/echo"
	"net/http"
	"io"
	"html/template"
	"github.com/anjunact/go-stock/models"
	"github.com/labstack/gommon/log"
	"path/filepath"
	"strconv"
	"fmt"
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

var templates map[string]*template.Template
func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	templatesDir := "./templates/"
	layouts, err := filepath.Glob(templatesDir + "layouts/*.html")
	if err != nil {
		log.Fatal(err)
	}
	widgets, err := filepath.Glob(templatesDir + "widgets/*.html")
	if err != nil {
		log.Fatal(err)
	}
	for _, layout := range layouts {
		files := append(widgets, layout)
		templates[filepath.Base(layout)] = template.Must(template.ParseFiles(files...))
	}
	fmt.Println(layouts)
}
func main() {
	rootPath,_:=  filepath.Abs(".")
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
		templates: template.Must(template.ParseGlob(rootPath+"/templates/*.html")),
	}
	e.Renderer = t
	e.GET("/stocks/:page", Stocks)

	e.Static("/static", rootPath+"/public")
	e.File("/",rootPath+"/public/index.html")
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