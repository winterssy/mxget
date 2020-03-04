package baidu

import (
	"context"
	"fmt"
	"strings"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
)

func (a *API) GetSong(ctx context.Context, songId string) (*api.Song, error) {
	resp, err := a.GetSongRaw(ctx, songId)
	if err != nil {
		return nil, err
	}

	resp.SongInfo.URL = songURL(resp.SongURL.URL)
	a.patchSongsLyric(ctx, &resp.SongInfo)
	songs := translate(&resp.SongInfo)
	return songs[0], nil
}

// 获取歌曲
func (a *API) GetSongRaw(ctx context.Context, songId string) (*SongResponse, error) {
	resp := new(SongResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiGetSong)
	req.SetQuery(aesCBCEncrypt(songId))
	req.SetContext(ctx)
	r, err := a.SendRequest(req)
	if err == nil {
		err = r.JSON(resp)
	}
	if err != nil {
		return nil, err
	}

	if resp.ErrorCode != 22000 {
		return nil, fmt.Errorf("get song: %s", resp.errorMessage())
	}

	return resp, nil
}

// 批量获取歌曲，遗留接口，不推荐使用
func (a *API) GetSongsRaw(ctx context.Context, songIds ...string) (*SongsResponse, error) {
	params := ghttp.Params{
		"songIds": strings.Join(songIds, ","),
	}
	resp := new(SongsResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiGetSongs)
	req.SetQuery(params)
	req.SetContext(ctx)
	r, err := a.SendRequest(req)
	if err == nil {
		err = r.JSON(resp)
	}
	if err != nil {
		return nil, err
	}

	if resp.ErrorCode != 22000 {
		return nil, fmt.Errorf("get songs: %d", resp.ErrorCode)
	}

	return resp, nil
}

func (a *API) GetSongLyric(ctx context.Context, songId string) (string, error) {
	resp, err := a.GetSongLyricRaw(ctx, songId)
	if err != nil {
		return "", err
	}

	return resp.LrcContent, nil
}

// 获取歌词
func (a *API) GetSongLyricRaw(ctx context.Context, songId string) (*SongLyricResponse, error) {
	params := ghttp.Params{
		"songid": songId,
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

	if resp.ErrorCode != 0 && resp.ErrorCode != 22000 {
		return nil, fmt.Errorf("get lyric: %s", resp.errorMessage())
	}

	return resp, nil
}
