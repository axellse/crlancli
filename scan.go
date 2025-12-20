package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v3"
)

func ScanForPrinters(con context.Context, cmd *cli.Command) error {
	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	masks := []net.IP{}
	ips := []net.IP{}
	ifaces := []string{}
	for _, i := range interfaces {
		if i.Flags & net.FlagLoopback != 0 {
			continue
		}

		a, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range a {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if ipNet.IP.To4() == nil {
					continue
				}
				masks = append(masks, net.IP(ipNet.Mask))
				ips = append(ips, ipNet.IP.To4())
				ifaces = append(ifaces, i.Name)
			}
		}
	}

	printers := []string{}
	startTime := time.Now()
	for i, rawIp := range ips {
		ipInt := binary.BigEndian.Uint32(rawIp)
		maskInt := binary.BigEndian.Uint32(masks[i])

		start := (ipInt & maskInt) +1
		end := (start | ^maskInt) -1

		for rip := start; rip < end; rip++ {
			ip := make(net.IP, 4)
			binary.BigEndian.PutUint32(ip, rip)

			go func() {
				socket, err := NewPrinterWebsocket(ip.String())
				if err != nil {
					return 
				}

				model, ok := Models[socket.InitalFrame.Model]
				if !ok {
					model = "\u001b[0;31mUnknown model! ID:" + socket.InitalFrame.Model + "\u001b[0m"
				}

				modelColor, ok := ModelColors[socket.InitalFrame.Model]
				if !ok {
					modelColor = "\u001b[0"
				}


				printers = append(printers, "Found " + modelColor + model + "\u001b[0m: \u001b[0;32m" + socket.InitalFrame.Hostname + "\u001b[0m (\u001b[0;32m" + ip.String() + "\u001b[0m)")			
				socket.Close()
			}()
		}
	}

	scanTime := 30
	if cmd.Bool("long") {
		scanTime = 150
	}
	for i := range scanTime {
		fmt.Print("\u001b[2K\r")
		fmt.Print("Scanning ", len(ifaces), " interface(s) for printers (" + strconv.FormatFloat(float64(scanTime - i) / 10, 'g', -1, 64) + "sec left)")
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Print("\u001b[2K\r")

	fmt.Println("Found", len(printers), "printers in", time.Since(startTime).Round(10 * time.Millisecond).Seconds(), "seconds:")
	fmt.Println(strings.Join(printers, "\n"))

	return nil
}