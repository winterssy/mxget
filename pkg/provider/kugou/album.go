package kugou

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
	albumInfo, err := a.GetAlbumInfoRaw(ctx, albumId)
	if err != nil {
		return nil, err
	}

	albumSongs, err := a.GetAlbumSongsRaw(ctx, albumId, 1, -1)
	if err != nil {
		return nil, err
	}

	n := len(albumSongs.Data.Info)
	if n == 0 {
		return nil, errors.New("get album songs: no data")
	}

	a.patchSongInfo(ctx, albumSongs.Data.Info...)
	a.patchSongsInfo(ctx, albumSongs.Data.Info...)
	a.patchSongsLyric(ctx, albumSongs.Data.Info...)
	songs := translate(albumSongs.Data.Info...)
	return &api.Collection{
		Id:     strconv.Itoa(albumInfo.Data.AlbumId),
		Name:   strings.TrimSpace(albumInfo.Data.AlbumName),
		PicURL: strings.ReplaceAll(albumInfo.Data.ImgURL, "{size}", "480"),
		Songs:  songs,
	}, nil
}

// 获取专辑信息
func (a *API) GetAlbumInfoRaw(ctx context.Context, albumId string) (*AlbumInfoResponse, error) {
	params := ghttp.Params{
		"albumid": albumId,
	}

	resp := new(AlbumInfoResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiGetAlbumInfo)
	req.SetQuery(params)
	req.SetContext(ctx)
	r, err := a.SendRequest(req)
	if err == nil {
		err = r.JSON(resp)
	}
	if err != nil {
		return nil, err
	}

	if resp.ErrCode != 0 {
		return nil, fmt.Errorf("get album info: %s", resp.errorMessage())
	}

	return resp, nil
}

// 获取专辑歌曲，page: 页码；pageSize: 每页数量，-1获取全部
func (a *API) GetAlbumSongsRaw(ctx context.Context, albumId string, page int, pageSize int) (*AlbumSongsResponse, error) {
	params := ghttp.Params{
		"albumid":  albumId,
		"page":     page,
		"pagesize": pageSize,
	}

	resp := new(AlbumSongsResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiGetAlbumSongs)
	req.SetQuery(params)
	req.SetContext(ctx)
	r, err := a.SendRequest(req)
	if err == nil {
		err = r.JSON(resp)
	}
	if err != nil {
		return nil, err
	}

	if resp.ErrCode != 0 {
		return nil, fmt.Errorf("get album songs: %s", resp.errorMessage())
	}

	return resp, nil
}
