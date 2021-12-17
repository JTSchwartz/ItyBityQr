package ityBityQr

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/golang/gddo/httputil/header"
	"github.com/skip2/go-qrcode"
)

type QrConfig struct {
	Url  string
	Size int
}

func ItyBityQr(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		config, err := ParseQuery(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		Respond(w, *config)

	case http.MethodPost:
		if r.Header.Get("Content-Type") != "" {
			value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
			if value != "application/json" {
				msg := "Content-Type header is not application/json"
				http.Error(w, msg, http.StatusUnsupportedMediaType)
				return
			}
		}

		config, err := ParseBody(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		Respond(w, *config)
	default:
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func Respond(w http.ResponseWriter, config QrConfig) {
	png, err := qrcode.Encode(config.Url, qrcode.Medium, config.Size)
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

func ParseBody(w http.ResponseWriter, r *http.Request) (*QrConfig, error) {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1048576))
	dec.DisallowUnknownFields()

	var config QrConfig
	if err := dec.Decode(&config); err != nil {
		HandleJSONParsingError(w, err)
		return nil, err
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		http.Error(w, msg, http.StatusBadRequest)
		return &config, err
	}

	return &config, nil
}

func ParseQuery(r *http.Request) (*QrConfig, error) {
	var config QrConfig
	query := r.URL.Query()

	url := query.Get("url")
	if url == "" {
		return &config, errors.New("Url Param 'url' is missing")
	}

	size, _ := strconv.Atoi(query.Get("size"))
	if size == 0 {
		return &config, errors.New("Url Param 'size' must be set")
	}

	config = BuildQrConfig(url, size)

	return &config, nil
}
