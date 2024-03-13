package main

import (
	"fmt"
	"jichinx/hik-sign/internal/sign"
)

const APP_SECRET = "DJMVuJhQjx1BABPyEmPa"
const (
	URL_CAMERA_SEARCH = "http://127.0.0.1/artemis/api/resource/v2/camera/search"
)

func main() {

	headers := map[string]string{
		sign.HEADER_X_CA_KEY:     "23752999",
		sign.HEADER_ACCEPT:       "*/*",
		sign.HEADER_CONTENT_TYPE: "application/json",
	}
	m, err := sign.ObtainSign("post", URL_CAMERA_SEARCH, APP_SECRET, headers)
	if err != nil {
		panic(err)
	}
	fmt.Println(m)
}
