package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/vooon/esp-ota-server/server"
	"os"
)

var opts struct {
	Bind    string `short:"s" long:"bind" env:"EOBIND" default:":8092" description:"bind address"`
	BaseUrl string `short:"u" long:"base-url" env:"EOBASEURL" default:"http://localhost:8092" description:"base url"`
	DataDir string `short:"d" long:"data-dir" env:"EODATADIR" required:"true" description:"path to data dir"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		if flerr, ok := err.(*flags.Error); ok && flerr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	}

	config := server.Config{
		Bind:        opts.Bind,
		BaseUrl:     opts.BaseUrl,
		DataDirPath: opts.DataDir,
	}

	server.Serve(config)
}
