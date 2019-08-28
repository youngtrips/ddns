package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/youngtrips/ddns/internal/config"
	"github.com/youngtrips/ddns/server"
	"github.com/youngtrips/ddns/version"
	"gopkg.in/urfave/cli.v2"
)

const (
	APP_NAME  = "ddns"
	APP_USAGE = "ddns"
)

func main() {
	app := &cli.App{
		Name:    APP_NAME,
		Usage:   APP_USAGE,
		Version: version.Version,

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: fmt.Sprintf("conf/%s.yml", APP_NAME),
				Usage: "specified config file",
			},
		},
		Action: func(c *cli.Context) error {
			if err := config.Load(c.String("config")); err != nil {
				log.Error("load config failed: ", err)
				return err
			}

			ctx, cancelFunc := context.WithCancel(context.TODO())
			sc := make(chan os.Signal, 1)
			signal.Notify(sc, os.Kill, syscall.SIGTERM, os.Interrupt, syscall.SIGINT, syscall.SIGHUP)
			log.Infof("start server [%s]", APP_NAME)
			go server.Run(ctx)
			select {
			case <-sc:
				cancelFunc()
				log.Infof("stop server [%s]", APP_NAME)
			}
			return nil
		},
	}
	app.Run(os.Args)
}
