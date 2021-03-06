package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/convox/cli/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/convox/cli/stdcli"
)

func init() {
	stdcli.RegisterCommand(cli.Command{
		Name:        "system",
		Description: "manage the base convox system",
		Usage:       "",
		Action:      cmdSystem,
		Subcommands: []cli.Command{
			{
				Name:        "update",
				Description: "update the convox system API",
				Usage:       "<version>",
				Action:      cmdSystemUpate,
			},
			{
				Name:        "scale",
				Description: "scale the convox system cluster",
				Usage:       "",
				Action:      cmdSystemScale,
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "count",
						Usage: "instance count, e.g. 3 or 10",
					},
					cli.StringFlag{
						Name:  "type",
						Usage: "instance type, e.g. t2.small or c3.xlargs",
					},
				},
			},
		},
	})
}

func cmdSystem(c *cli.Context) {
	data, err := ConvoxGet("/system")

	if err != nil {
		stdcli.Error(err)
		return
	}

	var a *App
	err = json.Unmarshal(data, &a)

	fmt.Printf("Name       %s\n", a.Name)
	fmt.Printf("Status     %s\n", a.Status)
	fmt.Printf("Version    %s\n", a.Parameters["Version"])
	fmt.Printf("Count      %s\n", a.Parameters["InstanceCount"])
	fmt.Printf("Type       %s\n", a.Parameters["InstanceType"])
}

func cmdSystemUpate(c *cli.Context) {
	version := ""

	if len(c.Args()) == 0 {
		resp, err := http.Get("http://convox.s3.amazonaws.com/release/latest/version")

		if err != nil {
			fmt.Printf("ERROR")
			stdcli.Error(err)
			return
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			stdcli.Error(err)
			return
		}

		version = string(body)
	} else {
		version = c.Args()[0]
	}

	v := url.Values{}
	v.Set("version", version)

	_, err := ConvoxPostForm("/system", v)

	if err != nil {
		stdcli.Error(err)
		return
	}

	cmdSystem(c)
}

func cmdSystemScale(c *cli.Context) {
	v := url.Values{}

	if c.IsSet("count") {
		v.Set("count", c.String("count"))
	}

	if c.IsSet("type") {
		v.Set("type", c.String("type"))
	}

	if len(v) > 0 {
		_, err := ConvoxPostForm("/system", v)

		if err != nil {
			stdcli.Error(err)
			return
		}
	}

	cmdSystem(c)
}
