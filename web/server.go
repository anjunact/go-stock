package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/anjunact/go-stock/models"
	"github.com/anjunact/go-stock/web/uitls"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
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
	// templatesDir := "./github.com/anjunact/go-stock/web/templates/"
	templatesDir := "./templates/"
	layouts, err := filepath.Glob(templatesDir + "layouts/index.html")
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
	fmt.Println(templates)

}
func main() {
	e := echo.New()
	// e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
	// 	return func(c echo.Context) error {
	// 		cc := &CustomContext{c}
	// 		return h(cc)
	// 	}
	// })
	// e.GET("/", func(c echo.Context) error {
	// 	cc := c.(*CustomContext)
	// 	cc.Foo()
	// 	cc.Bar()
	// 	return cc.String(200, "OK")
	// })
	templatesDir := "./templates/"
	layouts, err := filepath.Glob(templatesDir + "*.html")
	if err != nil {
		log.Fatal(err)
	}

	widgets, err := filepath.Glob(templatesDir + "widgets/*.html")
	if err != nil {
		log.Fatal(err)
	}
	files := append(widgets, layouts...)
	fmt.Println(files)
	t := &Template{
		templates: template.Must(template.ParseFiles(files...)),
		// templates: template.Must(template.ParseFiles(files...))，
	}
	e.Renderer = t
	e.GET("/stocks/:page", Stocks)

	e.Static("/static", "public")
	e.File("/", "public/index.html")
	e.Logger.Fatal(e.Start(":1323"))
}
func Stocks(c echo.Context) error {
	page := c.Param("page")
	pageInt, _ := strconv.Atoi(page)
	pageSize := 10
	stocks, err := models.Stocks(pageInt, pageSize)
	if err != nil {
		log.Print(err)
	}
	data := map[string]interface{}{}
	data["stocks"] = stocks

	stock := models.Stock{}
	count := stock.Count()
	fmt.Println(count)
	p := utils.NewPaginater(count, pageInt, pageSize)
	data["p"] = p
	fmt.Printf("%+v\n", p)

	return c.Render(http.StatusOK, "stocks", data)
}
