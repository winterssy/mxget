package xiami

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
	resp, err := a.GetAlbumRaw(ctx, albumId)
	if err != nil {
		return nil, err
	}

	_songs := resp.Data.Data.AlbumDetail.Songs
	n := len(_songs)
	if n == 0 {
		return nil, errors.New("get album: no data")
	}

	a.patchSongsLyric(ctx, _songs...)
	songs := translate(_songs...)
	return &api.Collection{
		Id:     resp.Data.Data.AlbumDetail.AlbumId,
		Name:   strings.TrimSpace(resp.Data.Data.AlbumDetail.AlbumName),
		PicURL: resp.Data.Data.AlbumDetail.AlbumLogo,
		Songs:  songs,
	}, nil
}

// 获取专辑
func (a *API) GetAlbumRaw(ctx context.Context, albumId string) (*AlbumResponse, error) {
	token, err := a.getToken(apiGetAlbum)
	if err != nil {
		return nil, err
	}

	model := make(map[string]string)
	_, err = strconv.Atoi(albumId)
	if err != nil {
		model["albumStringId"] = albumId
	} else {
		model["albumId"] = albumId
	}
	params := signPayload(token, model)

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

	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("get album: %s", err.Error())
	}

	return resp, nil
}
