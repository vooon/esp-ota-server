package main

import (
	"github.com/alecthomas/kong"

	"github.com/vooon/esp-ota-server/server"
)

var opts struct {
	Bind    string `short:"s" name:"bind" env:"EOBIND" default:":8092" help:"bind address"`
	BaseUrl string `short:"u" name:"base-url" env:"EOBASEURL" default:"http://localhost:8092" help:"base url"`
	DataDir string `short:"d" name:"data-dir" env:"EODATADIR" required:"true" help:"path to data dir"`
}

func main() {
	kctx := kong.Parse(&opts,
		kong.Description("ESP8266 & ESP32 Firmware OTA server"),
		kong.Configuration(kong.JSON, "/etc/espotad/espotad.json", "~/.config/espotad.json"),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: false,
		}),
		kong.DefaultEnvars("ESPOTA"),
	)

	config := server.Config{
		Bind:        opts.Bind,
		BaseUrl:     opts.BaseUrl,
		DataDirPath: opts.DataDir,
	}

	err := server.Serve(config)
	kctx.FatalIfErrorf(err)
}
