package main

import (
	"flag"
	"github.com/ian-kent/go-log/log"
	comcfg "github.com/jphautin/MailHog/config"
	"github.com/jphautin/mailhog-gui/web"
	"github.com/jphautin/mailhog-server/api"
	"github.com/jphautin/mailhog-server/config"
	"github.com/jphautin/mailhog-server/smtp"
	gohttp "net/http"
	"os"
)

var conf *config.Config
var comconf *comcfg.Config
var exitCh chan int

func configure() {
	comcfg.RegisterFlags()
	config.RegisterFlags()
	flag.Parse()
	conf = config.Configure()
	comconf = comcfg.Configure()
}

func main() {
	configure()

	if comconf.AuthFile != "" {
		web.AuthFile(comconf.AuthFile)
	}

	exitCh = make(chan int)
	cb := func(r gohttp.Handler) {
		api.CreateAPI(conf, r)
	}
	go web.Listen(conf.APIBindAddr, cb)
	go smtp.Listen(conf, exitCh)

	for {
		select {
		case <-exitCh:
			log.Printf("Received exit signal")
			os.Exit(0)
		}
	}
}
