package ityBityQr

import (
	"bytes"
	"errors"
	"image/color"
	"io"
	"net/http"
	"strconv"

	"github.com/skip2/go-qrcode"
)

type QrConfig struct {
	url  string
	size int
	bg   color.Color
	fg   color.Color
}

func ItyBityQr(w http.ResponseWriter, r *http.Request) {
	config, err, genErr := ParseQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if genErr != nil && *genErr == "" {
		w.Header().Set("GenerationError", *genErr)
	}

	png, err := qrcode.Encode(config.url, qrcode.Medium, config.size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusFailedDependency)
		return
	}

	_, err = io.Copy(w, bytes.NewBuffer(png))
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

	config = BuildQrConfig(url, size)

	return &config, nil, nil
}
