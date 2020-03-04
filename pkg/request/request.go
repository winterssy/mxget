package request

import (
	"github.com/winterssy/ghttp"
)

const (
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"
)

var (
	DefaultClient *ghttp.Client
)

func init() {
	DefaultClient = ghttp.New()
	DefaultClient.RegisterBeforeRequestCallbacks(
		ghttp.WithUserAgent(defaultUserAgent),
		ghttp.EnableRetry(),
	)
}
