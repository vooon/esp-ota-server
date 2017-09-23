package server

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/GeertJohan/go.rice"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
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
	lg := c.Logger()

	project := c.Param("project")
	filename := c.Param("file")

	path := filepath.Join(s.config.DataDirPath, project, filename)
	file, err := os.Open(path)
	if err != nil && os.IsNotExist(err) {
		lg.Warnj(log.JSON{
			"msg":       "File not found",
			"err":       err,
			"file_path": path,
		})
		return c.String(http.StatusNotFound, "no file")
	} else if err != nil {
		return err
	}
	defer file.Close()

	md5hasher := md5.New()
	sha512hasher := sha512.New()

	teeRd := io.TeeReader(io.TeeReader(file, md5hasher), sha512hasher)

	b, err := ioutil.ReadAll(teeRd)
	if err != nil {
		return err
	}

	md5sum := hex.EncodeToString(md5hasher.Sum(nil))
	sha512sum := hex.EncodeToString(sha512hasher.Sum(nil))

	hdr := c.Request().Header

	lg.Printj(log.JSON{
		"esp8266_request_headers": hdr,
	})

	//staMac, _ := hdr["X-Esp8266-Sta-Mac"]
	//apMac, _ := hdr["X-Esp8266-Ap-Mac"]
	//freeSpace, _ := hdr["X-Esp8266-Free-Space"]
	//sketchSize, _ := hdr["X-Esp8266-Sketch-Size"]
	sketchMd5, md5ok := hdr["X-Esp8266-Sketch-Md5"]
	//chipSize, _ := hdr["X-Esp8266-Chip-Size"]
	//sdkVersion, _ := hdr["X-Esp8266-Sdk-Version"]
	mode, ok := hdr["X-Esp8266-Mode"]
	version, vok := hdr["X-Esp8266-Version"]

	if !ok {
		return c.String(http.StatusBadRequest, "bad request")
	}

	sendFile := true
	if vok {
		vmap := map[string]string{}
		for _, kv := range strings.Split(version[0], " ") {
			n := strings.SplitN(kv, ":", 2)
			vmap[n[0]] = n[1]
		}

		c.Logger().Printj(log.JSON{
			"esp8266_version_map": vmap,
		})

		// if version has MD5
		md5, mok := vmap["md5"]
		if mok {
			sendFile = md5 != md5sum
		}
	}
	if md5ok {
		sendFile = sketchMd5[0] != md5sum
	}

	c.Response().Header()["x-MD5"] = []string{md5sum} // do not do strings.Title()
	c.Response().Header().Set("x-SHA512", sha512sum)  // not used by actual version
	lg.Printj(log.JSON{
		"esp8266_mode": mode[0],
		"send_file":    sendFile,
		"file_path":    path,
		"file_size":    len(b),
	})

	if sendFile {
		//return c.Blob(http.StatusOK, "application/ocetet-stream", b)
		return c.File(path)
	} else {
		return c.String(http.StatusNotModified, "")
	}
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

	newpath, err := filepath.Abs(config.DataDirPath)
	if err != nil {
		e.Logger.Fatal("can't abs data-dir")
	}
	if stat, err := os.Stat(newpath); err == nil && stat.IsDir() {
		e.Logger.Info("Data-dir: ", newpath)
		config.DataDirPath = newpath
	} else {
		e.Logger.Fatal("data-dir not exist! ", newpath)
	}

	var templates *template.Template = nil
	assets := rice.MustFindBox("../assets")
	assets.Walk("", func(name string, info os.FileInfo, err error) error {
		e.Logger.Info("Processing asset ", name)

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
