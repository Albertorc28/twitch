package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/kbinani/screenshot"
)

//Config es la configuración del cliente
type Config struct {
	Protocol        string
	Host            string
	Port            string
	Streamer        string
	FrameIntervalMs time.Duration
}

func main() {
	n := screenshot.NumActiveDisplays()
	fmt.Printf("Pantallas: %d", n)
	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	var config Config

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		panic(err)
	}

	for {
		inicio := time.Now()
		time.Sleep(config.FrameIntervalMs * time.Millisecond)
		//bounds := screenshot.GetDisplayBounds(0)
		img, err := screenshot.CaptureDisplay(0)
		if err != nil {
			panic(err)
		}

		err = postFrame(img, &config)

		if err == nil {
			elapsed := time.Since(inicio)
			fmt.Printf("Frame enviado (%s)\n", elapsed)
		} else {
			fmt.Println("fail")
		}

	}
}

func postFrame(img *image.RGBA, config *Config) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	formWriter, err := bodyWriter.CreateFormFile("frame", "frame.jpg")

	jpeg.Encode(formWriter, img, nil)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	targetURL := config.Protocol + "://" + config.Host + ":" + config.Port + "/stream/" + config.Streamer + "/addframe"
	resp, err := http.Post(targetURL, contentType, bodyBuf)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}
