package main

import (
	"encoding/json"
	"fmt"
	"jichinx/hik-sign/internal/consts"
	"jichinx/hik-sign/internal/entity/body"
	"jichinx/hik-sign/internal/hikhttp"
	"sync"
	"time"
)

type Obj map[string]interface{}

func main() {
	// 1，先查询根目录
	fmt.Println("1, 先查询根目录")
	rootIndexCode, err := fistGetRoot()
	Must(err)

	// 2, 查询所
	fmt.Println("2, 查询区域信息")
	region2s, err := getRegionsByKey([]string{rootIndexCode.(string)}, nil)
	Must(err)
	region2Codes := filterRegionCodes(region2s)
	fmt.Println(region2Codes)

	// 3，查询所下面的 农贸市场、食品流通类型
	fmt.Println("3, 查询所下面的 农贸市场、食品流通类型")
	region3s, err := getRegionsByKey(region2Codes, []string{"农贸市场", "食品流通"})
	Must(err)
	region3Codes := filterRegionCodes(region3s)
	fmt.Println(region3Codes)

	// 4, 查询具体场所
	fmt.Println("4, 查询具体场所")
	region4s, err := getRegionsByKey(region3Codes, nil)
	Must(err)
	region4Codes := filterRegionCodes(region4s)
	fmt.Println(region4Codes)

	// 5, 查询场所下的摄像头
	fmt.Println("5, 查询场所下的摄像头")
	cameras, err := getCameras(region4Codes)
	Must(err)
	cameraCodes := filterCameraCodes(cameras)
	fmt.Println(len(cameraCodes))
	fmt.Println(cameraCodes)
	// 6, 取预览流地址
	fmt.Println("6, 取预览流地址")
	previewUrls, err := getCameraPreviewUrls(cameraCodes)
	Must(err)
	fmt.Println(previewUrls...)
}

func filterRegionCodes(regions []interface{}) []string {
	indexCodes := make([]string, 0)
	if regions == nil {
		return indexCodes
	}
	// codeNames := make([]string, 0)
	for _, v := range regions {
		regionObj := v.(map[string]interface{})
		indexCodes = append(indexCodes, regionObj["indexCode"].(string))
		// codeNames = append(codeNames, regionObj["name"].(string))
	}
	// fmt.Println(codeNames)
	return indexCodes
}

func filterCameraCodes(cameras []interface{}) []string {
	indexCodes := make([]string, 0)
	if cameras == nil {
		return indexCodes
	}
	// cameraNames := make([]string, 0)
	for _, v := range cameras {
		regionObj := v.(map[string]interface{})
		indexCodes = append(indexCodes, regionObj["indexCode"].(string))
		// cameraNames = append(cameraNames, regionObj["name"].(string))
	}
	// fmt.Println(cameraNames)
	return indexCodes
}

func getCameras(parentIndexCodes []string) ([]interface{}, error) {
	body := body.Values{}
	body.Add("pageNo", 1)
	body.Add("pageSize", 1000)
	body.Add("regionIndexCodes", parentIndexCodes)

	//请求数据
	respBytes, err := hikhttp.Request(createUrl(consts.SERVER_HOST, consts.URL_CAMERA_SEARCH), body.String())
	if err != nil {
		return nil, err
	}
	//显示数据
	var respObj Obj
	err = json.Unmarshal(respBytes, &respObj)
	if err != nil {
		return nil, err
	}
	// fmt.Println("---", respObj)
	regions := respObj["data"].(map[string]interface{})
	// fmt.Println("---", regions)
	regionList := regions["list"].([]interface{})

	return regionList, nil
}
func getCameraPreviewUrls(cameraCodes []string) ([]interface{}, error) {

	urls := make([]interface{}, 0)

	requestPreviewUrl := func(code string) error {
		body := body.Values{}
		body.Add("pageNo", 1)
		body.Add("pageSize", 1000)
		body.Add("cameraIndexCode", code)
		//请求数据
		respBytes, err := hikhttp.Request(createUrl(consts.SERVER_HOST, consts.URL_CAMERA_PREVIEW), body.String())
		if err != nil {
			return err
		}
		//显示数据
		var respObj Obj
		err = json.Unmarshal(respBytes, &respObj)
		if err != nil {
			return err
		}
		fmt.Println("---", respObj)
		regions := respObj["data"].(map[string]interface{})
		// fmt.Println("---", regions)
		previewUrl := regions["url"].(string)
		urls = append(urls, previewUrl)
		return nil
	}

	timeNow := time.Now()
	var wg sync.WaitGroup
	wg.Add(len(cameraCodes))
	for _, v := range cameraCodes {
		go func() {
			requestPreviewUrl(v)
			wg.Done()
		}()
	}
	wg.Wait()
	s := time.Since(timeNow)
	fmt.Println(s)

	return urls, nil
}

func getRegionsByKey(parentIndexCodes []string, keyWords []string) ([]interface{}, error) {
	// fmt.Println(parentIndexCodes, keyWords)
	list := make([]interface{}, 0)

	doRequest := func(body body.Values) ([]interface{}, error) {
		// fmt.Println(body.String())
		//请求数据
		respBytes, err := hikhttp.Request(createUrl(consts.SERVER_HOST, consts.URL_REGION_SEARCH), body.String())
		if err != nil {
			return nil, err
		}
		//显示数据
		var respObj Obj
		err = json.Unmarshal(respBytes, &respObj)
		if err != nil {
			return nil, err
		}
		// fmt.Println("---", respObj)
		regions := respObj["data"].(map[string]interface{})
		// fmt.Println("---", regions)
		regionList := regions["list"].([]interface{})
		// fmt.Println("---", regionList)
		// for _, v := range regionList {
		// regionObj := v.(map[string]interface{})
		// fmt.Println("-----", regionObj)
		// }
		list = append(list, regionList...)
		// rootIndexCode := regionRoot["indexCode"]
		// fmt.Println(regions)
		// 根据根目录 查询区域信息
		// return rootIndexCode, nil
		return list, nil
	}

	body := body.Values{}
	body.Add("pageNo", 1)
	body.Add("pageSize", 1000)
	body.Add("parentIndexCodes", parentIndexCodes)
	body.Add("resourceType", "camera")
	if keyWords == nil {
		doRequest(body)
	} else {
		for _, keyWord := range keyWords {
			body.Add("regionName", keyWord)
			doRequest(body)
		}
	}
	return list, nil
}

func fistGetRoot() (interface{}, error) {
	body := body.Values{}
	body.Add("pageNo", 1)
	body.Add("pageSize", 10)
	//请求数据
	respBytes, err := hikhttp.Request(createUrl(consts.SERVER_HOST, consts.URL_REGION_ROOT), body.String())
	if err != nil {
		return nil, err
	}
	//显示数据
	var respObj Obj
	err = json.Unmarshal(respBytes, &respObj)
	if err != nil {
		return nil, err
	}
	regionRoot := respObj["data"].(map[string]interface{})
	fmt.Println(regionRoot)
	rootIndexCode := regionRoot["indexCode"]
	fmt.Println(rootIndexCode)
	// 根据根目录 查询区域信息
	return rootIndexCode, nil
}
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
func createUrl(host string, path string) string {
	return fmt.Sprintf("%s/artemis/%s", host, path)
}
