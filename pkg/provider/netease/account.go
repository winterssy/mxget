package netease

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/winterssy/ghttp"
)

// 邮箱登录
func (a *API) EmailLoginRaw(ctx context.Context, email string, password string) (*LoginResponse, error) {
	passwordHash := md5.Sum([]byte(password))
	password = hex.EncodeToString(passwordHash[:])
	data := map[string]interface{}{
		"username":      email,
		"password":      password,
		"rememberLogin": true,
	}

	resp := new(LoginResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodPost, apiEmailLogin)
	req.SetForm(weapi(data))
	req.SetContext(ctx)
	r, err := a.SendRequest(req)
	if err == nil {
		err = r.JSON(resp)
	}
	if err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("email login: %s", resp.errorMessage())
	}

	return resp, nil
}

// 手机登录
func (a *API) CellphoneLoginRaw(ctx context.Context, countryCode int, phone int, password string) (*LoginResponse, error) {
	passwordHash := md5.Sum([]byte(password))
	password = hex.EncodeToString(passwordHash[:])
	data := map[string]interface{}{
		"phone":         phone,
		"countrycode":   countryCode,
		"password":      password,
		"rememberLogin": true,
	}

	resp := new(LoginResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodPost, apiCellphoneLogin)
	req.SetForm(weapi(data))
	req.SetContext(ctx)
	r, err := a.SendRequest(req)
	if err == nil {
		err = r.JSON(resp)
	}
	if err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("cellphone login: %s", resp.errorMessage())
	}

	return resp, nil
}

// 刷新登录状态
func (a *API) RefreshLoginRaw(ctx context.Context) (*CommonResponse, error) {
	resp := new(CommonResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodPost, apiRefreshLogin)
	req.SetForm(weapi(struct{}{}))
	req.SetContext(ctx)
	r, err := a.SendRequest(req)
	if err == nil {
		err = r.JSON(resp)
	}
	if err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("refresh login: %s", resp.errorMessage())
	}

	return resp, nil
}

// 退出登录
func (a *API) LogoutRaw(ctx context.Context) (*CommonResponse, error) {
	resp := new(CommonResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodPost, apiLogout)
	req.SetForm(weapi(struct{}{}))
	req.SetContext(ctx)
	r, err := a.SendRequest(req)
	if err == nil {
		err = r.JSON(resp)
	}
	if err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("logout: %s", resp.errorMessage())
	}

	return resp, nil
}
