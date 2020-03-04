package baidu

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

	n := len(resp.Result.SongList)
	if n == 0 {
		return nil, errors.New("get playlist: no data")
	}

	a.patchSongsURL(ctx, resp.Result.SongList...)
	a.patchSongsLyric(ctx, resp.Result.SongList...)
	songs := translate(resp.Result.SongList...)
	return &api.Collection{
		Id:     resp.Result.Info.ListId,
		Name:   strings.TrimSpace(resp.Result.Info.ListTitle),
		PicURL: resp.Result.Info.ListPic,
		Songs:  songs,
	}, nil
}

// 获取歌单
func (a *API) GetPlaylistRaw(ctx context.Context, playlistId string) (*PlaylistResponse, error) {
	params := ghttp.Params{
		"list_id":   playlistId,
		"withcount": "1",
		"withsong":  "1",
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

	if resp.ErrorCode != 22000 {
		return nil, fmt.Errorf("get playlist: %s", resp.errorMessage())
	}

	return resp, nil
}
