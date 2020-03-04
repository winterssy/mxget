package netease

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
)

func (a *API) GetAlbum(ctx context.Context, albumId string) (*api.Collection, error) {
	_albumId, err := strconv.Atoi(albumId)
	if err != nil {
		return nil, err
	}

	resp, err := a.GetAlbumRaw(ctx, _albumId)
	if err != nil {
		return nil, err
	}

	n := len(resp.Songs)
	if n == 0 {
		return nil, errors.New("get album: no data")
	}

	a.patchSongsURL(ctx, songDefaultBR, resp.Songs...)
	a.patchSongsLyric(ctx, resp.Songs...)
	songs := translate(resp.Songs...)
	return &api.Collection{
		Id:     strconv.Itoa(resp.Album.Id),
		Name:   strings.TrimSpace(resp.Album.Name),
		PicURL: resp.Album.PicURL,
		Songs:  songs,
	}, nil
}

// 获取专辑
func (a *API) GetAlbumRaw(ctx context.Context, albumId int) (*AlbumResponse, error) {
	resp := new(AlbumResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodPost, fmt.Sprintf(apiGetAlbum, albumId))
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
		return nil, fmt.Errorf("get album: %s", resp.errorMessage())
	}

	return resp, nil
}
