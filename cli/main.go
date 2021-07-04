package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/alecthomas/kong"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

type Context struct {
	Debug      bool
	DryRun     bool
	IpServer   string
	ZoneId     string
	RecordName string
}

func GetIp(ctx *Context) (string, error) {
	resp, err := http.Get(ctx.IpServer)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

type OnceCmd struct {
	Arg string
}

func (o *OnceCmd) Run(ctx *Context) error {
	fmt.Println("run command!")
	fmt.Println(ctx.IpServer)

	ip, err := GetIp(ctx)
	if err != nil {
		return err
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("us-west-2"),
		},
	})
	if err != nil {
		return err
	}
	svc := route53.New(sess)
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(ctx.RecordName),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(ip),
							},
						},
						TTL:  aws.Int64(60),
						Type: aws.String("A"),
					},
				},
			},
			Comment: aws.String("DDNS update"),
		},
		HostedZoneId: aws.String(ctx.ZoneId),
	}

	result, err := svc.ChangeResourceRecordSets(input)
	if err != nil {
		return err
	}

	fmt.Println(result)

	return nil
}

type ContinuousCmd struct{}

func (c *ContinuousCmd) Run(ctx *Context) error {
	fmt.Println("run continuously")
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
