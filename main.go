package main

import (
	"fmt"
	"jichinx/hik-sign/internal/consts"
	"jichinx/hik-sign/internal/entity/body"
	"jichinx/hik-sign/internal/hikhttp"
)

func main() {
	//先查询根目录
	body := body.Values{}
	body.Add("pageNo", 1)
	body.Add("pageSize", 10)
	regionRoot, err := hikhttp.Request(createUrl(consts.SERVER_HOST, consts.URL_REGION_ROOT), body.String())
	Must(err)
	fmt.Println(regionRoot)

}
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
func createUrl(host string, path string) string {
	return fmt.Sprintf("%s/artemis/%s", host, path)
}
