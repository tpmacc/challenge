package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/tpmacc/challenge"
	"html/template"
	"io"
	"net/http"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if err := challenge.RenderTemplate(w, name, data); err != nil {
		c.Logger().Error(err)
	}

	return nil
//	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Renderer = &Template{templates:challenge.Templates}

	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "yoyo.html", challenge.ReadData())
	})
	e.GET("/league", func(c echo.Context) error {
		return c.Render(http.StatusOK, "leagueTable.html", challenge.ReadData())
	})
	e.GET("/wager", func(c echo.Context) error {
		return c.Render(http.StatusOK, "wager.html", challenge.BuildChallengeStandings())
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}


