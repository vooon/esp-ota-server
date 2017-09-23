package server

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/GeertJohan/go.rice"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//go:generate rice embed-go

type server struct {
	config    Config
	templates *template.Template
}

// Render renders a template document
func (s server) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return s.templates.ExecuteTemplate(w, name, data)
}

func (s server) getBinaryFile(c echo.Context) error {
	return c.String(http.StatusOK, "Hello world")
}

func (s server) get403(c echo.Context) error {
	return c.Render(http.StatusForbidden, "403.ghtm", map[string]interface{}{
		"BarbradyJpgBase64": "/assets/barbrady.jpg",
	})
}

func Serve(config Config) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	var templates *template.Template = nil
	assets := rice.MustFindBox("../assets")
	assets.Walk("/", func(name string, info os.FileInfo, err error) error {
		e.Logger.Debug("Processing asset ", name)

		if m, _ := filepath.Match("*.ghtm", name); !m {
			return nil
		}

		bn := filepath.Base(name)
		s, _ := assets.String(name)

		var tmpl *template.Template
		if templates == nil {
			templates = template.New(bn)
		}

		if bn == templates.Name() {
			tmpl = templates
		} else {
			tmpl = templates.New(bn)
		}

		_, err2 := tmpl.Parse(s)
		return err2
	})

	s := server{
		config:    config,
		templates: templates,
	}

	assetHandler := http.FileServer(assets.HTTPBox())

	e.Renderer = s
	e.GET("/bin/:project/:file", s.getBinaryFile)
	//e.POST("/bin/:project/:file", postBinaryFile)
	e.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/assets/", assetHandler)))
	e.GET("/", s.get403)

	e.Logger.Fatal(e.Start(config.Bind))
}
