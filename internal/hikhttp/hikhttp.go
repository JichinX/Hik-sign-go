package hikhttp

import (
	"io"
	"jichinx/hik-sign/internal/consts"
	"jichinx/hik-sign/internal/sign"
	"net/http"
	"strings"
)

func Request(url string, body string) ([]byte, error) {

	var respBytes []byte
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return respBytes, err
	}
	err = patchHeaders(req)
	if err != nil {
		return respBytes, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return respBytes, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)

}
func patchHeaders(req *http.Request) (err error) {
	signHeaders, err := sign.ObtainSign("POST", req.URL.String(), consts.APP_SECRET, consts.CommHeaders)
	if err != nil {
		return
	}
	for k, v := range consts.CommHeaders {
		// fmt.Println("add header", k, v)
		req.Header.Add(k, v)
	}
	for k, v := range signHeaders {
		// fmt.Println("add header", k, v)
		req.Header.Add(k, v)
	}
	return
}
