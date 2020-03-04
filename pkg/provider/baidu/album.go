package baidu

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
)

func (a *API) GetAlbum(ctx context.Context, albumId string) (*api.Collection, error) {
	resp, err := a.GetAlbumRaw(ctx, albumId)
	if err != nil {
		return nil, err
	}

	n := len(resp.SongList)
	if n == 0 {
		return nil, errors.New("get album: no data")
	}

	a.patchSongsURL(ctx, resp.SongList...)
	a.patchSongsLyric(ctx, resp.SongList...)
	songs := translate(resp.SongList...)
	return &api.Collection{
		Id:     resp.AlbumInfo.AlbumId,
		Name:   strings.TrimSpace(resp.AlbumInfo.Title),
		PicURL: strings.SplitN(resp.AlbumInfo.PicBig, "@", 2)[0],
		Songs:  songs,
	}, nil
}

// 获取专辑
func (a *API) GetAlbumRaw(ctx context.Context, albumId string) (*AlbumResponse, error) {
	params := ghttp.Params{
		"album_id": albumId,
	}

	resp := new(AlbumResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiGetAlbum)
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
		return nil, fmt.Errorf("get album: %s", resp.errorMessage())
	}

	return resp, nil
}
