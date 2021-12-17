package ityBityQr

func BuildQrConfig(url string, size int) QrConfig {
	return QrConfig{url, size}
}
