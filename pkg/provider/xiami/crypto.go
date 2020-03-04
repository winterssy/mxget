package xiami

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/gjson"
)

const (
	APPKey = "23649156"
)

var (
	reqHeader = map[string]interface{}{
		"appId":      200,
		"platformId": "h5",
	}
)

func signPayload(token string, model interface{}) ghttp.Params {
	payload := map[string]interface{}{
		"header": reqHeader,
		"model":  model,
	}
	requestStr, _ := gjson.MarshalToString(payload)
	data := map[string]string{
		"requestStr": requestStr,
	}
	dataStr, _ := gjson.MarshalToString(data)

	t := time.Now().UnixNano() / (1e6)
	signStr := fmt.Sprintf("%s&%d&%s&%s", token, t, APPKey, dataStr)
	sign := fmt.Sprintf("%x", md5.Sum([]byte(signStr)))

	return ghttp.Params{
		"t":    t,
		"sign": sign,
		"data": dataStr,
	}
}
