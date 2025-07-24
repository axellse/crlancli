package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v3"
)

func SendAction(con context.Context, cmd *cli.Command) error {
	fileba, err := os.ReadFile(cmd.StringArg("file"))
	if err == os.ErrNotExist {
		fmt.Println("error: file provided does not exist.")
		return err
	} else if err != nil {
		fmt.Println("error: could not read provided file.")
		return err
	}

	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	name := filepath.Base(cmd.StringArg("file"))

	fw, err := form.CreateFormFile("file", name)
	if err != nil {
		fmt.Println("error: could not add file to form.")
		return err
	}

	_, werr := fw.Write(fileba)
	if werr != nil {
		fmt.Println("error: could not write file to form.")
		return err
	}

	cerr := form.Close()
	if cerr != nil {
		fmt.Println("error: could not close to form.")
		return err
	}

	transferreq, err := http.NewRequest("POST", "http://"+cmd.StringArg("printer")+"/upload/"+name, &body)
	if err != nil {
		fmt.Println("error: could not create request")
		return err
	}
	transferreq.Header.Add("Content-Type", form.FormDataContentType())

	resp, err := http.DefaultClient.Do(transferreq)
	if err != nil {
		fmt.Println("error: could not send request:", err)
		return err
	}

	if resp.StatusCode != 200 {
		ba, _ := io.ReadAll(resp.Body)
		fmt.Println("error: printer reported http error:", string(ba))
		return err
	}

	fmt.Println("file has been transferred to the printer with the name", name)

	if cmd.Bool("print") {
		time.Sleep(3 * time.Second)
		pws := NewPrinterWebsocket(cmd.StringArg("printer"))
		err := pws.ModifyFile(name, "print")
		if err != nil {
			log.Fatal("failed printing uploaded file:", err)
		}
		fmt.Println("print has been started")
	}
	return nil
}
