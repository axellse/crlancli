package main

import (
	"encoding/json"
	"log"
	"strings"

	"golang.org/x/net/websocket"
)

type printerSocketRequest struct {
	Method string `json:"method"`
	Params map[string]any `json:"params"`
}

type PrinterSocket struct {
	ws *websocket.Conn
}


//sends a raw query to the printer
func (ps PrinterSocket) SendRawQuery(method string, params map[string]any) error {
	payload := printerSocketRequest{
		Method: method,
		Params: params,
	}

	ba, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, werr := ps.ws.Write(ba)
	if werr != nil {
		return werr
	}

	return nil
}

//modifies a gcode file (delete, rename, print)
func (ps PrinterSocket) ModifyFile(file string, action string) error {
	return ps.SendRawQuery("set", map[string]any{
		"opGcodeFile" : strings.ToLower(action) + "prt:/usr/data/printer_data/gcodes/" + file,
	})
}

func NewPrinterWebsocket(url string) PrinterSocket {
	ws, err := websocket.Dial("ws://" + url + ":9999", "", "http://" + url)
	if err != nil {
		log.Fatal("could not create printer websocket, is the printer online?", err)
	}

	return PrinterSocket{
		ws: ws,
	}
}