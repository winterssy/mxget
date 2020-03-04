package netease

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/gjson"
	"github.com/winterssy/mxget/pkg/api"
)

func (a *API) GetSong(ctx context.Context, songId string) (*api.Song, error) {
	_songId, err := strconv.Atoi(songId)
	if err != nil {
		return nil, err
	}

	resp, err := a.GetSongsRaw(ctx, _songId)
	if err != nil {
		return nil, err
	}
	if len(resp.Songs) == 0 {
		return nil, errors.New("get song: no data")
	}

	_song := resp.Songs[0]
	a.patchSongsURL(ctx, songDefaultBR, _song)
	a.patchSongsLyric(ctx, _song)
	songs := translate(_song)
	return songs[0], nil
}

// 批量获取歌曲详情，上限1000首
func (a *API) GetSongsRaw(ctx context.Context, songIds ...int) (*SongsResponse, error) {
	n := len(songIds)
	if n > songRequestLimit {
		songIds = songIds[:songRequestLimit]
		n = songRequestLimit
	}

	c := make([]map[string]int, n)
	for i, id := range songIds {
		c[i] = map[string]int{"id": id}
	}
	enc, _ := gjson.MarshalToString(c)
	data := map[string]string{
		"c": enc,
	}

	resp := new(SongsResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodPost, apiGetSongs)
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
		return nil, fmt.Errorf("get songs: %s", resp.errorMessage())
	}

	return resp, nil
}

func (a *API) GetSongURL(ctx context.Context, id int, br int) (string, error) {
	resp, err := a.GetSongsURLRaw(ctx, br, id)
	if err != nil {
		return "", err
	}
	if len(resp.Data) == 0 {
		return "", errors.New("get song url: no data")
	}
	if resp.Data[0].Code != 200 {
		return "", errors.New("get song url: copyright protection")
	}

	return resp.Data[0].URL, nil
}

// 批量获取歌曲播放地址，br: 比特率，128/192/320/999
func (a *API) GetSongsURLRaw(ctx context.Context, br int, songIds ...int) (*SongURLResponse, error) {
	var _br int
	switch br {
	case 128, 192, 320:
		_br = br
	default:
		_br = 999
	}

	enc, _ := gjson.MarshalToString(songIds)
	data := map[string]interface{}{
		"br":  _br * 1000,
		"ids": enc,
	}

	resp := new(SongURLResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodPost, apiGetSongsURL)
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
		return nil, fmt.Errorf("get songs url: %s", resp.errorMessage())
	}

	return resp, nil
}

func (a *API) GetSongLyric(ctx context.Context, songId int) (string, error) {
	resp, err := a.GetSongLyricRaw(ctx, songId)
	if err != nil {
		return "", err
	}
	return resp.Lrc.Lyric, nil
}

// 获取歌词
func (a *API) GetSongLyricRaw(ctx context.Context, songId int) (*SongLyricResponse, error) {
	data := map[string]interface{}{
		"method": "POST",
		"url":    "https://music.163.com/api/song/lyric?lv=-1&kv=-1&tv=-1",
		"params": map[string]int{
			"id": songId,
		},
	}
	resp := new(SongLyricResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodPost, linuxAPI)
	req.SetForm(linuxapi(data))
	req.SetContext(ctx)
	r, err := a.SendRequest(req)
	if err == nil {
		err = r.JSON(resp)
	}
	if err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("get song lyric: %s", resp.errorMessage())
	}

	return resp, nil
}
