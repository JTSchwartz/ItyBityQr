package ityBityQr

import (
	"image/color"
)

func BuildQrConfig(url string, size int, colors ...color.Color) QrConfig {
	if len(colors) != 2 {
		return QrConfig{url, size, color.White, color.Black}
	}
	return QrConfig{url, size, colors[0], colors[1]}
}
