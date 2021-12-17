package ityBityQr

import (
	"fmt"
	"image/color"
)

func ParseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 6:
		_, err = fmt.Sscanf(s, "%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
	case 3:
		_, err = fmt.Sscanf(s, "%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")

	}
	return
}

func BuildQrConfig(url string, size int, colors ...color.Color) QrConfig {
	if len(colors) != 2 {
		return QrConfig{url, size, color.White, color.Black}
	}
	return QrConfig{url, size, colors[0], colors[1]}
}

func ValidateUniqueColors(set color.Color, check color.Color, backup color.Color) color.Color {
	if set == check {
		return backup
	}
	return check
}
