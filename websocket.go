package main

import (
	"encoding/json"
	"errors"
	"strings"

	"golang.org/x/net/websocket"
)

type printerSocketRequest struct {
	Method string `json:"method"`
	Params map[string]any `json:"params"`
}

type PrinterSocket struct {
	ws *websocket.Conn
	InitalFrame InitalWSResponse
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

//closes the socket
func (ps PrinterSocket) Close() error {
	return ps.ws.Close()
}

//modifies a gcode file (delete, rename, print)
func (ps PrinterSocket) ModifyFile(file string, action string) error {
	return ps.SendRawQuery("set", map[string]any{
		"opGcodeFile" : strings.ToLower(action) + "prt:/usr/data/printer_data/gcodes/" + file,
	})
}

func NewPrinterWebsocket(url string) (PrinterSocket, error) {
	ws, err := websocket.Dial("ws://" + url + ":9999", "", "http://" + url)
	if err != nil {
		return PrinterSocket{}, errors.New("could not create printer websocket, is the printer online?\n" + err.Error())
	}

	var InitalData InitalWSResponse
	err = websocket.JSON.Receive(ws, &InitalData)
	if err != nil {
		return PrinterSocket{}, errors.New("failed receiving initial frame from websocket, what?\n" + err.Error())
	}

	return PrinterSocket{
		ws: ws,
		InitalFrame: InitalData,
	}, nil
}