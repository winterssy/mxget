package migu

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

	if len(resp.Resource) == 0 || len(resp.Resource[0].SongItems) == 0 {
		return nil, errors.New("get album: no data")
	}

	a.patchSongsLyric(ctx, resp.Resource[0].SongItems...)
	songs := translate(resp.Resource[0].SongItems...)
	return &api.Collection{
		Id:     resp.Resource[0].AlbumId,
		Name:   strings.TrimSpace(resp.Resource[0].Title),
		PicURL: picURL(resp.Resource[0].ImgItems),
		Songs:  songs,
	}, nil
}

// 获取专辑
func (a *API) GetAlbumRaw(ctx context.Context, albumId string) (*AlbumResponse, error) {
	params := ghttp.Params{
		"resourceId": albumId,
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

	if resp.Code != "000000" {
		return nil, fmt.Errorf("get album: %s", resp.errorMessage())
	}

	return resp, nil
}
