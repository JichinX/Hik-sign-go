package consts

import "jichinx/hik-sign/internal/sign"

const (
	APP_SECRET  = "DJMVuJhQjx1BABPyEmPa"
	APP_KEY     = "23752999"
	SERVER_HOST = "http://112.6.118.61:8081"
)

const (
	URL_CAMERA_SEARCH = "api/resource/v2/camera/search" //查询监控点信息
	URL_REGION_ROOT   = "api/resource/v1/regions/root"  //区域根

	URL_REGION_SEARCH  = "api/irds/v2/region/nodesByParams" //区域信息查询
	URL_CAMERA_PREVIEW = "api/video/v2/cameras/previewURLs" //获取预览地址
)

var CommHeaders = map[string]string{
	sign.HEADER_X_CA_KEY:     APP_KEY,
	sign.HEADER_ACCEPT:       "*/*",
	sign.HEADER_CONTENT_TYPE: "application/json",
}
