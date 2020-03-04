package tencent

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/gjson"
	"github.com/winterssy/mxget/pkg/api"
)

func (a *API) GetSong(ctx context.Context, songMid string) (*api.Song, error) {
	resp, err := a.GetSongRaw(ctx, songMid)
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, errors.New("get song: no data")
	}

	_song := resp.Data[0]
	a.patchSongsURLV1(ctx, _song)
	a.patchSongsLyric(ctx, _song)
	songs := translate(_song)
	return songs[0], nil
}

// 获取歌曲详情
func (a *API) GetSongRaw(ctx context.Context, songMid string) (*SongResponse, error) {
	params := ghttp.Params{
		"songmid": songMid,
	}

	resp := new(SongResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiGetSong)
	req.SetQuery(params)
	req.SetContext(ctx)
	r, err := a.SendRequest(req)
	if err == nil {
		err = r.JSON(resp)
	}
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("get song: %d", resp.Code)
	}

	return resp, nil
}

func (a *API) GetSongURLV1(ctx context.Context, songMid string, mediaMid string) (string, error) {
	const (
		tmpl = "http://mobileoc.music.tc.qq.com/%s?guid=0&uin=0&vkey=%s"
	)

	resp, err := a.GetSongURLV1Raw(ctx, songMid, mediaMid)
	if err != nil {
		return "", err
	}
	if len(resp.Data.Items) == 0 {
		return "", errors.New("get song url: no data")
	}

	item := resp.Data.Items[0]
	if item.SubCode != 0 {
		return "", fmt.Errorf("get song url: %d", item.SubCode)
	}

	return fmt.Sprintf(tmpl, item.FileName, item.Vkey), nil
}

// 获取歌曲播放地址
func (a *API) GetSongURLV1Raw(ctx context.Context, songMid string, mediaMid string) (*SongURLResponseV1, error) {
	params := ghttp.Params{
		"songmid":  songMid,
		"filename": "M500" + mediaMid + ".mp3",
	}

	resp := new(SongURLResponseV1)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiGetSongURLV1)
	req.SetQuery(params)
	req.SetContext(ctx)
	r, err := a.SendRequest(req)
	if err == nil {
		err = r.JSON(resp)
	}
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("get song url: %s", resp.ErrInfo)
	}

	return resp, nil
}

func (a *API) GetSongURLV2(ctx context.Context, songMid string) (string, error) {
	resp, err := a.GetSongsURLV2Raw(ctx, songMid)
	if err != nil {
		return "", err
	}
	if len(resp.Req0.Data.MidURLInfo) == 0 {
		return "", errors.New("get song url: no data")
	}

	n := len(resp.Req0.Data.Sip)
	if n == 0 {
		return "", errors.New("get song url: no sip")
	}

	// 随机获取一个sip
	sip := resp.Req0.Data.Sip[rand.Intn(n)]
	urlInfo := resp.Req0.Data.MidURLInfo[0]
	if urlInfo.PURL == "" {
		return "", errors.New("get song url: copyright protection")
	}

	return sip + urlInfo.PURL, nil
}

// 批量获取歌曲播放地址
func (a *API) GetSongsURLV2Raw(ctx context.Context, songMids ...string) (*SongURLResponseV2, error) {
	if len(songMids) > songURLRequestLimit {
		songMids = songMids[:songURLRequestLimit]
	}

	param := map[string]interface{}{
		"guid":      "9386820628",
		"loginflag": 1,
		"songmid":   songMids,
		"uin":       "1152921504838414858",
		"platform":  "20",
	}
	req0 := map[string]interface{}{
		"module": "vkey.GetVkeyServer",
		"method": "CgiGetVkey",
		"param":  param,
	}
	data := map[string]interface{}{
		"req0": req0,
	}

	enc, _ := gjson.MarshalToString(data)
	params := ghttp.Params{
		"data": enc,
	}
	resp := new(SongURLResponseV2)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiGetSongsURLV2)
	req.SetQuery(params)
	req.SetContext(ctx)
	r, err := a.SendRequest(req)
	if err == nil {
		err = r.JSON(resp)
	}
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("get song url: %d", resp.Code)
	}

	return resp, nil
}

func (a *API) GetSongLyric(ctx context.Context, songMid string) (string, error) {
	resp, err := a.GetSongLyricRaw(ctx, songMid)
	if err != nil {
		return "", err
	}

	// lyric, err := base64.StdEncoding.DecodeString(resp.Lyric)
	// if err != nil {
	// 	return "", err
	// }

	return resp.Lyric, nil
}

// 获取歌词
func (a *API) GetSongLyricRaw(ctx context.Context, songMid string) (*SongLyricResponse, error) {
	params := ghttp.Params{
		"songmid": songMid,
	}

	resp := new(SongLyricResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiGetSongLyric)
	req.SetQuery(params)
	req.SetContext(ctx)
	r, err := a.SendRequest(req)
	if err == nil {
		err = r.JSON(resp)
	}
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("get song lyric: %d", resp.Code)
	}

	return resp, nil
}
