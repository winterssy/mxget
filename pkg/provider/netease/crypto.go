package netease

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/gjson"
	"github.com/winterssy/mxget/pkg/cryptography"
)

const (
	Base62                      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	PresetKey                   = "0CoJUm6Qyw8W8jud"
	IV                          = "0102030405060708"
	LinuxAPIKey                 = "rFgB&h#%2?^eDg:Q"
	EAPIKey                     = "e82ckenh8dichen8"
	DefaultRSAPublicKeyModulus  = "e0b509f6259df8642dbc35662901477df22677ec152b5ff68ace615bb7b725152b3ab17a876aea8a5aa76d2e417629ec4ee341f56135fccf695280104e0312ecbda92557c93870114af6c9d05c4f7f0c3685b7a46bee255932575cce10b424d813cfe4875d3e82047b97ddef52741d546b8e289dc6935b3ece0462db0a22b8e7"
	DefaultRSAPublicKeyExponent = 0x10001
)

func weapi(origData interface{}) ghttp.Form {
	plainText, _ := gjson.Marshal(origData)
	params := base64.StdEncoding.EncodeToString(cryptography.AESCBCEncrypt(plainText, []byte(PresetKey), []byte(IV)))
	secKey := CreateSecretKey(16, Base62)
	params = base64.StdEncoding.EncodeToString(cryptography.AESCBCEncrypt([]byte(params), secKey, []byte(IV)))
	return ghttp.Form{
		"params":    params,
		"encSecKey": cryptography.RSAEncrypt(BytesReverse(secKey), DefaultRSAPublicKeyModulus, DefaultRSAPublicKeyExponent),
	}
}

func linuxapi(origData interface{}) ghttp.Form {
	plainText, _ := gjson.Marshal(origData)
	return ghttp.Form{
		"eparams": strings.ToUpper(hex.EncodeToString(cryptography.AESECBEncrypt(plainText, []byte(LinuxAPIKey)))),
	}
}

func eapi(url string, origData interface{}) ghttp.Form {
	plainText, _ := gjson.MarshalToString(origData)
	message := fmt.Sprintf("nobody%suse%smd5forencrypt", url, plainText)
	digest := fmt.Sprintf("%x", md5.Sum([]byte(message)))
	data := fmt.Sprintf("%s-36cd479b6b5-%s-36cd479b6b5-%s", url, plainText, digest)
	return ghttp.Form{
		"params": strings.ToUpper(hex.EncodeToString(cryptography.AESECBEncrypt([]byte(data), []byte(EAPIKey)))),
	}
}

func CreateSecretKey(size int, charset string) []byte {
	secKey, n := make([]byte, size), len(charset)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range secKey {
		secKey[i] = charset[r.Intn(n)]
	}
	return secKey
}

func BytesReverse(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}
