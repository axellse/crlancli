package main

import (
	"context"
	"log"
	"os"
	"github.com/urfave/cli/v3"
)

func main() {
    cmd := &cli.Command{
		Name: "crlancli",
		Description: "a cli utility to interact with networked creality printers",
		Usage: "interact with networked creality printers",
		Commands: []*cli.Command{
			{
				Name:  "scan",
				Usage: "scans the network for any printers",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name: "long",
						Aliases: []string{"l"},
						Usage: "scans for 15 seconds instead of just 3",
					},
				},
				Action: ScanForPrinters,
			},
			{
				Name:  "temp",
				Usage: "temprature readings for a printer",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name: "printer",
						UsageText: "[printer ip address/hostname, eg. 192.168.1.218] ",
					},
				},
				Action: TempAction,
			},
			{
				Name:  "send",
				Usage: "transfer a gcode file to a printer",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name: "file",
						UsageText: "[gcode file] ",
					},
					&cli.StringArg{
						Name: "printer",
						UsageText: "[printer ip address/hostname, eg. 192.168.1.218] ",
					},
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name: "print",
						Aliases: []string{"p"},
						Usage: "whether or not to start printing the file after transferring",
					},
				},
				Action: SendAction,
			},
		},
	} 

    if err := cmd.Run(context.Background(), os.Args); err != nil {
        log.Fatal(err)
    }
}