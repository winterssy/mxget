package xiami

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/winterssy/ghttp"
)

// 登录接口，account 可为邮箱/手机号码
func (a *API) LoginRaw(ctx context.Context, account string, password string) (*LoginResponse, error) {
	token, err := a.getToken(apiLogin)
	if err != nil {
		return nil, err
	}

	passwordHash := md5.Sum([]byte(password))
	password = hex.EncodeToString(passwordHash[:])
	model := map[string]string{
		"account":  account,
		"password": password,
	}
	params := signPayload(token, model)

	resp := new(LoginResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiLogin)
	req.SetQuery(params)
	req.SetContext(ctx)
	r, err := a.SendRequest(req)
	if err == nil {
		err = r.JSON(resp)
	}
	if err != nil {
		return nil, err
	}

	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("login: %s", err.Error())
	}

	return resp, nil
}
