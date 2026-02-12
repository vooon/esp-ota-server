package main

import (
	"github.com/alecthomas/kong"
	"github.com/prometheus/common/version"

	"github.com/vooon/esp-ota-server/server"
)

var opts struct {
	Version    kong.VersionFlag `name:"version" help:"Show version and exit."`
	Bind       string `short:"s" name:"bind" env:"EOBIND" default:":8092" help:"bind address"`
	BaseUrl    string `short:"u" name:"base-url" env:"EOBASEURL" default:"http://localhost:8092" help:"base url"`
	DataDir    string `short:"d" name:"data-dir" env:"EODATADIR" required:"true" help:"path to data dir"`
	Prometheus bool   `short:"m" name:"prometheus" env:"EOPROMETHEUS" negatable:"" help:"Enable/disable prometheus /metrics endpoint"`
}

func main() {
	kctx := kong.Parse(&opts,
		kong.Description("ESP8266 & ESP32 Firmware OTA server"),
		kong.Configuration(kong.JSON, "/etc/espotad/espotad.json", "~/.config/espotad.json"),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: false,
		}),
		kong.DefaultEnvars("ESPOTA"),
		kong.Vars{
			"version": version.Print("espotad"),
		},
	)

	config := server.Config{
		Bind:             opts.Bind,
		BaseUrl:          opts.BaseUrl,
		DataDirPath:      opts.DataDir,
		EnablePrometheus: opts.Prometheus,
	}

	err := server.Serve(config)
	kctx.FatalIfErrorf(err)
}
