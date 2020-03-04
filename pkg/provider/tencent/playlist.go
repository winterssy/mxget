package tencent

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
)

func (a *API) GetPlaylist(ctx context.Context, playlistId string) (*api.Collection, error) {
	resp, err := a.GetPlaylistRaw(ctx, playlistId)
	if err != nil {
		return nil, err
	}
	if len(resp.Data.CDList) == 0 || len(resp.Data.CDList[0].SongList) == 0 {
		return nil, errors.New("get playlist: no data")
	}

	playlist := resp.Data.CDList[0]
	if playlist.PicURL == "" {
		playlist.PicURL = playlist.Logo
	}
	_songs := playlist.SongList
	a.patchSongsURLV1(ctx, _songs...)
	a.patchSongsLyric(ctx, _songs...)
	songs := translate(_songs...)
	return &api.Collection{
		Id:     playlist.DissTid,
		Name:   strings.TrimSpace(playlist.DissName),
		PicURL: playlist.PicURL,
		Songs:  songs,
	}, nil
}

// 获取歌单
func (a *API) GetPlaylistRaw(ctx context.Context, playlistId string) (*PlaylistResponse, error) {
	params := ghttp.Params{
		"id": playlistId,
	}

	resp := new(PlaylistResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiGetPlaylist)
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
		return nil, fmt.Errorf("get playlist: %d", resp.Code)
	}

	return resp, nil
}
