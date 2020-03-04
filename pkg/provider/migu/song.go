package migu

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
)

/*
	注意！
	GetSongIdRaw, GetSongPicRaw, GetSongLyricRaw 为网页版API
	这些API限流，并发请求经常503，不适用于批量获取
*/

func (a *API) GetSong(ctx context.Context, id string) (*api.Song, error) {
	/*
		先判断id是版权id还是歌曲id，以减少1次API请求
		测试表现版权id的长度是11位，以6开头并且可能包含字符，歌曲id为纯数字，长度不定
		不确定是否会误判，待反馈
	*/
	var songId string
	var err error
	if len(id) > 10 && strings.HasPrefix(id, "6") {
		songId, err = a.GetSongId(ctx, id)
		if err != nil {
			return nil, err
		}
	} else {
		songId = id
	}

	resp, err := a.GetSongRaw(ctx, songId)
	if err != nil {
		return nil, err
	}
	if len(resp.Resource) == 0 {
		return nil, errors.New("get song: no data")
	}

	_song := resp.Resource[0]

	// 单曲请求可调用网页版API获取歌词，不会出现乱码现象
	lyric, err := a.GetSongLyric(ctx, _song.CopyrightId)
	if err == nil {
		_song.Lyric = lyric
	}
	songs := translate(_song)
	return songs[0], nil
}

// 获取歌曲详情
func (a *API) GetSongRaw(ctx context.Context, songId string) (*SongResponse, error) {
	params := ghttp.Params{
		"songId": songId,
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

	if resp.Code != "000000" {
		return nil, fmt.Errorf("get song: %s", resp.errorMessage())
	}

	return resp, nil
}

func (a *API) GetSongId(ctx context.Context, copyrightId string) (string, error) {
	resp, err := a.GetSongIdRaw(ctx, copyrightId)
	if err != nil {
		return "", err
	}
	if len(resp.Items) == 0 || resp.Items[0].SongId == "" {
		return "", errors.New("get song id: no data")
	}

	return resp.Items[0].SongId, nil
}

// 根据版权id获取歌曲id
func (a *API) GetSongIdRaw(ctx context.Context, copyrightId string) (*SongIdResponse, error) {
	params := ghttp.Params{
		"copyrightId": copyrightId,
	}

	resp := new(SongIdResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiGetSongId)
	req.SetQuery(params)
	req.SetContext(ctx)
	r, err := a.SendRequest(req)
	if err == nil {
		err = r.JSON(resp)
	}
	if err != nil {
		return nil, err
	}

	if resp.ReturnCode != "000000" {
		return nil, fmt.Errorf("get song id: %s", resp.Msg)
	}

	return resp, nil
}

func (a *API) GetSongURL(ctx context.Context, contentId, resourceType string) (string, error) {
	resp, err := a.GetSongURLRaw(ctx, contentId, resourceType)
	if err != nil {
		return "", err
	}

	return resp.Data.URL, nil
}

// 获取歌曲播放地址
func (a *API) GetSongURLRaw(ctx context.Context, contentId, resourceType string) (*SongURLResponse, error) {
	params := ghttp.Params{
		"contentId":             contentId,
		"lowerQualityContentId": contentId,
		"resourceType":          resourceType,
	}

	resp := new(SongURLResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiGetSongURL)
	req.SetQuery(params)
	req.SetContext(ctx)
	r, err := a.SendRequest(req)
	if err == nil {
		err = r.JSON(resp)
	}
	if err != nil {
		return nil, err
	}

	if resp.Code != "000000" {
		return nil, fmt.Errorf("get song url: %s", resp.errorMessage())
	}

	return resp, nil
}

func (a *API) GetSongPic(ctx context.Context, songId string) (string, error) {
	resp, err := a.GetSongPicRaw(ctx, songId)
	if err != nil {
		return "", err
	}
	return resp.LargePic, nil
}

// 获取歌曲专辑封面
func (a *API) GetSongPicRaw(ctx context.Context, songId string) (*SongPicResponse, error) {
	params := ghttp.Params{
		"songId": songId,
	}

	resp := new(SongPicResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiGetSongPic)
	req.SetQuery(params)
	req.SetContext(ctx)
	r, err := a.SendRequest(req)
	if err == nil {
		err = r.JSON(resp)
	}
	if err != nil {
		return nil, err
	}

	if resp.ReturnCode != "000000" {
		return nil, fmt.Errorf("get song pic: %s", resp.Msg)
	}

	return resp, nil
}

func (a *API) GetSongLyric(ctx context.Context, copyrightId string) (string, error) {
	resp, err := a.GetSongLyricRaw(ctx, copyrightId)
	if err != nil {
		return "", err
	}
	return resp.Lyric, nil
}

// 获取歌词
func (a *API) GetSongLyricRaw(ctx context.Context, copyrightId string) (*SongLyricResponse, error) {
	params := ghttp.Params{
		"copyrightId": copyrightId,
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

	if resp.ReturnCode != "000000" {
		return nil, fmt.Errorf("get song lyric: %s", resp.Msg)
	}

	return resp, nil
}
