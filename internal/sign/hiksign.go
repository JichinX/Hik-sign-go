package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

const (
	HEADER_X_CA_KEY = "x-ca-key"
)

// 官方指定的不参与签名计算的 header
const (
	HEADER_X_CA_SIGN             = "X-Ca-Signature"
	HEADER_X_CA_SIGN_HEADERS     = "X-Ca-Signature-Headers"
	HEADER_ACCEPT                = "Accept"
	HEADER_CONTENT_MD5           = "Content-MD5"
	HEADER_CONTENT_TYPE          = "Content-Type"
	HEADER_DATE                  = "Date"
	HEADER_CONTENT_LEN           = "Content-Length"
	HEADER_SERVER                = "Server"
	HEADER_CONNECTION            = "Connection"
	HEADER_HOST                  = "Host"
	HEADER_TRANSFER_ENCODING     = "Transfer-Encoding"
	HEADER_X_APPLICATION_CONTENT = "X-Application-Context"
	HEADER_CONTENT_ENCODING      = "Content-Encoding"
)

var noNeededHeaders = [13]string{
	HEADER_DATE,
	HEADER_HOST,
	HEADER_ACCEPT,
	HEADER_SERVER,
	HEADER_X_CA_SIGN,
	HEADER_CONNECTION,
	HEADER_CONTENT_LEN,
	HEADER_CONTENT_MD5,
	HEADER_CONTENT_TYPE,
	HEADER_CONTENT_ENCODING,
	HEADER_TRANSFER_ENCODING,
	HEADER_X_CA_SIGN_HEADERS,
	HEADER_X_APPLICATION_CONTENT,
}

// ObtainSign
// method
//
// remoteUrl
//
// appSecret
//
// headers
func ObtainSign(method string, remoteUrl string, appSecret string, headers map[string]string) (map[string]string, error) {
	signMap := make(map[string]string, 1)
	elememts := make([]string, 0)
	//从请求头获取 appkey
	appKey := headers[HEADER_X_CA_KEY]
	if appKey == "" {
		return signMap, fmt.Errorf("headers not define %s", HEADER_X_CA_KEY)
	}
	// Http Method
	elememts = append(elememts, strings.TrimSpace(strings.ToUpper(method)))
	// 定义一个函数，概括重复逻辑
	tryHeader := func(key string) {
		accept, ok := headers[key]
		if ok {
			elememts = append(elememts, accept)
			delete(headers, key)
		}
	}
	// Accept、Content-MD5、Content-Type、Date
	// 有 key就取值，为空也无所谓，无值就不添加
	tryHeader(HEADER_ACCEPT)
	tryHeader(HEADER_CONTENT_MD5)
	tryHeader(HEADER_CONTENT_TYPE)
	tryHeader(HEADER_DATE)
	// 先将参与 Headers 签名计算的 Header 的 Key 转换为小写字母，
	// 然后按照字典排序后使用如下方式拼接:
	//    如果某个 Header 的 value 为空，则使用HeaderKey+“:”+“\n”参与签名，需要保留 Key 和英文冒号。
	//    需要去除value字符串头尾的空字符串，如 key=“a-key”,value="　abc　“，则参与签名时的字符串为"a-key:abc”。
	//    注意，参与headers签名计算的header的key在签名字符串中必须转换为小写字母。
	sortedKeys, newHeaders := prepareHeaders(headers)
	for _, v := range sortedKeys {
		elememts = append(elememts, fmt.Sprintf("%s:%s", v, newHeaders[v]))
	}
	// URL 处理
	// url 指 Path + Query + BodyForm 中 Form 参数，
	// 组织方法：对 Query+BodyForm 参数按照字典对 Key 进行排序后按照如下方法拼接，
	//         如果 Query 或 BodyForm 参数为空，则 Url = Path，不需要添加”?”，
	//         如果某个参数的 Value 为空只保留 Key 参与签名。
	u, err := url.Parse(remoteUrl)
	if err != nil {
		return signMap, fmt.Errorf("Parse remoteUrl:%s failed,\n %w", remoteUrl, err)
	}
	fmt.Println("u.Path", u.Path)
	fmt.Println("u.RawQuery", u.RawQuery)
	var uPath string
	if u.RawQuery != "" {
		uPath = fmt.Sprintf("%s?%s", u.Path, u.RawQuery)
	} else {
		uPath = u.Path
	}
	elememts = append(elememts, uPath)
	// 计算签名
	// //拼接字符串
	message := strings.Join(elememts, "\n")

	fmt.Println("Message:================start")
	fmt.Println(message)
	fmt.Println("Message:================end")
	// //生成
	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write([]byte(message))
	result := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	signMap[HEADER_X_CA_SIGN] = result

	return signMap, nil
}

// 去除不参数计算签名的 header
// header key 转为小写
func prepareHeaders(src map[string]string) ([]string, map[string]string) {
	newHeaders := make(map[string]string)
	sortedKeys := make([]string, 0)
	for k, v := range src {
		if i := arrayHas(noNeededHeaders[:], k); i != -1 {
			continue
		}
		// 转小写
		// 去空格
		newV := strings.TrimSpace(strings.ToLower(v))
		newHeaders[k] = newV
		sortedKeys = append(sortedKeys, k)
	}
	// 排序
	sort.Strings(sortedKeys)
	return sortedKeys, newHeaders
}

// 判断是否存在某个元素
func arrayHas(slice []string, find string) int {
	for i := range slice {
		if slice[i] == find {
			return i
		}
	}
	return -1
}
