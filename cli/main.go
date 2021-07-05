package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/featheredtoast/ddns-route53/internal/aws"
	"github.com/featheredtoast/ddns-route53/internal/iplookup"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Context struct {
	Debug      bool
	DryRun     bool
	IpServer   string
	ZoneId     string
	RecordName string
}

type OnceCmd struct {
}

func (cmd *OnceCmd) Run(ctx *Context) error {
	fmt.Println("run command!")
	fmt.Println(ctx.IpServer)

	i := &iplookup.IpGetter{Server: ctx.IpServer, Record: ctx.RecordName}
	updateNeeded, ip, err := i.IpChanged()
	if err != nil {
		return err
	}
	if !updateNeeded {
		fmt.Println("No update needed. returning.")
		return nil
	}
	updater := aws.IpUpdater{Ip: ip, RecordName: ctx.RecordName, ZoneId: ctx.ZoneId}
	result, err := updater.UpdateIp()
	if err != nil {
		return err
	}

	fmt.Println(result)

	return nil
}

type ContinuousCmd struct{}

func (cmd *ContinuousCmd) Run(ctx *Context) error {
	fmt.Println("run continuously")
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	timeout := 1 * time.Second
runner:
	for {
		select {
		case <-time.After(timeout):
			fmt.Println("checking IP")
			checker := OnceCmd{}
			err := checker.Run(ctx)
			if err != nil {
				fmt.Println("Got error: ", err)
			}
			timeout = 10 * time.Minute
		case <-sig:
			fmt.Println("exiting")
			break runner
		}
	}
	return nil
}

var cli struct {
	Debug      bool          `help:"enable debug"`
	DryRun     bool          `help:"print out new IP info, but do not update."`
	IpServer   string        `help:"IP address lookup server. Defaults to https://ipinfo.io/ip" default:"https://ipinfo.io/ip"`
	ZoneId     string        `help:"Route 53 zone to update" required:"" env:"DDNS_R53_ZONE_ID"`
	RecordName string        `help:"The record name to update" required:"" env:"DDNS_R53_RECORD_NAME"`
	Once       OnceCmd       `cmd help:"Run once" default:"1"`
	Continuous ContinuousCmd `cmd help:"Run continuously"`
}

func main() {
	ctx := kong.Parse(&cli)
	err := ctx.Run(&Context{Debug: cli.Debug, DryRun: cli.DryRun, IpServer: cli.IpServer, ZoneId: cli.ZoneId, RecordName: cli.RecordName})
	ctx.FatalIfErrorf(err)
}
