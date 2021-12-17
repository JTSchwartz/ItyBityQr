package ityBityQr

import (
	"bytes"
	"errors"
	"image/color"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/skip2/go-qrcode"
)

type QrConfig struct {
	url  string
	size int
	bg   color.Color
	fg   color.Color
}

func QrGenerator(w http.ResponseWriter, r *http.Request) {
	config, err, genErr := ParseQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if *genErr == "" {
		w.Header().Set("GenerationError", *genErr)
	}

	err = qrcode.WriteColorFile(config.url, qrcode.Medium, config.size, config.bg, config.fg, "qr.png")
	if err != nil {
		http.Error(w, err.Error(), http.StatusFailedDependency)
		return
	}

	encoded, err := StreamImageFile()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(w, encoded)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ParseQuery(r *http.Request) (*QrConfig, error, *string) {
	var config QrConfig
	query := r.URL.Query()

	url := query.Get("url")
	if url == "" {
		return &config, errors.New("Url Param 'url' is missing"), nil
	}

	size, _ := strconv.Atoi(query.Get("size"))
	if size == 0 {
		return &config, errors.New("Url Param 'size' must be set"), nil
	}

	bg := query.Get("bg")
	fg := query.Get("fg")

	if bg != "" && fg != "" {
		bg, fg, genErr := ParseColors(bg, fg)
		config := BuildQrConfig(url, size, bg, fg)

		return &config, nil, &genErr
	}

	config = BuildQrConfig(url, size)

	return &config, nil, nil
}

func ParseColors(bgHex string, fgHex string) (bgColor color.Color, fgColor color.Color, err string) {
	bgColor, bgErr := ParseHexColor(bgHex)
	fgColor, fgErr := ParseHexColor(fgHex)

	if bgErr != nil && fgErr != nil {
		err = "Unable to parse background and foreground colors"
		bgColor = color.White
		fgColor = color.Black
	} else if bgErr != nil {
		err = "Unable to parse background color"
		bgColor = ValidateUniqueColors(fgColor, color.White, color.Black)
	} else if fgErr != nil {
		err = "Unable to parse background color"
		fgColor = ValidateUniqueColors(bgColor, color.Black, color.White)
	} else if bgColor == fgColor {
		if fgColor == color.Black {
			fgColor = color.White
		}
		fgColor = color.Black
	}

	return
}

func StreamImageFile() (stream *bytes.Buffer, err error) {
	file, err := os.Open("qr.png")
	if err != nil {
		return
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		return
	}

	buffer := make([]byte, fileinfo.Size())

	_, err = file.Read(buffer)
	if err != nil {
		return
	}

	stream = bytes.NewBuffer(buffer)
	return
}
