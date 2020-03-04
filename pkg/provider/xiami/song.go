package xiami

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
)

func (a *API) GetSong(ctx context.Context, songId string) (*api.Song, error) {
	resp, err := a.GetSongDetailRaw(ctx, songId)
	if err != nil {
		return nil, err
	}

	_song := &resp.Data.Data.SongDetail
	a.patchSongsLyric(ctx, _song)
	songs := translate(_song)
	return songs[0], nil
}

// 获取歌曲详情
func (a *API) GetSongDetailRaw(ctx context.Context, songId string) (*SongDetailResponse, error) {
	token, err := a.getToken(apiGetSongDetail)
	if err != nil {
		return nil, err
	}

	model := make(map[string]string)
	_, err = strconv.Atoi(songId)
	if err != nil {
		model["songStringId"] = songId
	} else {
		model["songId"] = songId
	}
	params := signPayload(token, model)

	resp := new(SongDetailResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiGetSongDetail)
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
		return nil, fmt.Errorf("get song detail: %s", err.Error())
	}

	return resp, nil
}

func (a *API) GetSongLyric(ctx context.Context, songId string) (string, error) {
	resp, err := a.GetSongLyricRaw(ctx, songId)
	if err != nil {
		return "", err
	}

	for _, i := range resp.Data.Data.Lyrics {
		if i.FlagOfficial == "1" && i.Type == "2" {
			return i.Content, nil
		}
	}

	return "", errors.New("official lyric not present")
}

// 获取歌词
func (a *API) GetSongLyricRaw(ctx context.Context, songId string) (*SongLyricResponse, error) {
	token, err := a.getToken(apiGetSongLyric)
	if err != nil {
		panic(err)
	}

	model := make(map[string]string)
	_, err = strconv.Atoi(songId)
	if err != nil {
		model["songStringId"] = songId
	} else {
		model["songId"] = songId
	}
	params := signPayload(token, model)

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

	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("get song lyric: %s", err.Error())
	}

	return resp, nil
}

// 批量获取歌曲，上限200首
func (a *API) GetSongsRaw(ctx context.Context, songIds ...string) (*SongsResponse, error) {
	token, err := a.getToken(apiGetSongs)
	if err != nil {
		return nil, err
	}

	if len(songIds) > songRequestLimit {
		songIds = songIds[:songRequestLimit]
	}
	model := map[string][]string{
		"songIds": songIds,
	}
	params := signPayload(token, model)

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

	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("get songs: %s", err.Error())
	}

	return resp, nil
}
