package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/urfave/cli/v3"
)

func roundStr(input string) string {
	f, err := strconv.ParseFloat(input, 64)
	if err != nil {
		fmt.Println("fatal: printer sending invalid temp readings")
		os.Exit(1)
	}
	return strconv.FormatFloat(f, 'g', -1, 64)
}

type Thermostat struct {
	Name string
	Reading string
	Target string
	Max string
}

func TempAction(con context.Context, cmd *cli.Command) error {
	pws, err := NewPrinterWebsocket(cmd.StringArg("printer"))
	if err != nil {
		fmt.Println("error: printer reported ws error:", err)
		return err
	}

	tstats := []Thermostat{
		{
			Name: "Bed 1: ",
			Reading: pws.InitalFrame.BedTemp0,
			Target: strconv.Itoa(pws.InitalFrame.TargetBedTemp0),
			Max: strconv.Itoa(pws.InitalFrame.MaxBedTemp),
		},
		{
			Name: "Bed 2: ",
			Reading: pws.InitalFrame.BedTemp1,
			Target: strconv.Itoa(pws.InitalFrame.TargetBedTemp1),
			Max: strconv.Itoa(pws.InitalFrame.MaxBedTemp),
		},
		{
			Name: "Bed 3: ",
			Reading: pws.InitalFrame.BedTemp2,
			Target: strconv.Itoa(pws.InitalFrame.TargetBedTemp2),
			Max: strconv.Itoa(pws.InitalFrame.MaxBedTemp),
		},
		{
			Name: "Nozzle:",
			Reading: pws.InitalFrame.NozzleTemp,
			Target: strconv.Itoa(pws.InitalFrame.TargetNozzleTemp),
			Max: strconv.Itoa(pws.InitalFrame.MaxNozzleTemp),
		},
		{
			Name: "Box:   ",
			Reading: strconv.Itoa(pws.InitalFrame.BoxTemp),
			Target: "N/A",
			Max: "N/A",
		},
	}

	for _, t := range tstats {
		if roundStr(t.Reading) == "0" {continue}
		fmt.Println(t.Name + "\u001b[0;32m", roundStr(t.Reading) + "°", "\u001b[0m→\u001b[0;31m", roundStr(t.Target) + "°", "\u001b[0m(Max " + roundStr(t.Max) + "°)")
	}

	return pws.Close()
}
