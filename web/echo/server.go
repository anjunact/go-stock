package main

import (
	"io"
	"net/http"
	"os"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"time"
)

// e.GET("/users/:id", getUser)
func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

//e.GET("/show", show)
//show?team=x-men&member=wolverine
func show(c echo.Context) error {
	// Get team and member from the query string
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team:"+team+", member:"+member)
}

// e.POST("/save", save)
//curl -F "name=Joe Smith" -F "email=joe@labstack.com" http://localhost:1323/save

func save(c echo.Context) error {
	// Get name and email
	name := c.FormValue("name")
	email := c.FormValue("email")
	return c.String(http.StatusOK, "name:"+name+", email:"+email)
}

//curl -F "name=Joe Smith" -F "avatar=@/Users/anjun/imgs/0.png" http://localhost:1323/save2
func save2(c echo.Context) error {
	// Get name
	name := c.FormValue("name")
	// Get avatar
	avatar, err := c.FormFile("avatar")
	if err != nil {
		return err
	}

	// Source
	src, err := avatar.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(avatar.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "<b>Thank you! "+name+"</b>")
}

type User struct {
	Name  string `json:"name" xml:"name" form:"name" query:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"name"`
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/users", func(c echo.Context) error {
		u := new(User)
		if err := c.Bind(u); err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, u)
		// or
		// return c.XML(http.StatusCreated, u)
	})
	// e.POST("/users", saveUser)
	// e.GET("/users/:id", getUser)
	// e.PUT("/users/:id", updateUser)
	// e.DELETE("/users/:id", deleteUser)
	e.POST("/save", save)
	e.POST("/save2", save2)
	e.GET("/users/:id", getUser)
	e.GET("/show", show)
	e.Static("/static", "static")

	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Group level middleware
	g := e.Group("/admin")
	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) ( bool) {
		if username == "joe" && password == "secret" {
			return  true
		}
		return false
	}))

	// Route level middleware
	track := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			println("request to /users")
			return next(c)
		}
	}
	e.GET("/users", func(c echo.Context) error {
		return c.String(http.StatusOK, "/users")
	}, track)
	//curl -X POST http://localhost:1323/users -H 'Content-Type: application/json' -d '{"name":"Joe","email":"joe@labstack"}'
	e.POST("/users",func(c echo.Context) (err error) {
		u := new(User)
		if err = c.Bind(u); err != nil {
			return
		}
		return c.JSON(http.StatusOK, u)
	})

	//l, err := net.Listen("tcp", ":1323")
	//if err != nil {
	//	e.Logger.Fatal(l)
	//}
	//e.Listener = l
	//e.Logger.Fatal(e.Start(""))
	s := &http.Server{
		Addr:         ":1324",
		ReadTimeout:  20 * time.Minute,
		WriteTimeout: 20 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))
	//e.Logger.Fatal(e.Start(":1324"))
}
