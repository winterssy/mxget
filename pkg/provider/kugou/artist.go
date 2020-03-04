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

func (a *API) GetArtist(ctx context.Context, singerId string) (*api.Collection, error) {
	artistInfo, err := a.GetArtistInfoRaw(ctx, singerId)
	if err != nil {
		return nil, err
	}

	artistSongs, err := a.GetArtistSongsRaw(ctx, singerId, 1, 50)
	if err != nil {
		return nil, err
	}

	n := len(artistSongs.Data.Info)
	if n == 0 {
		return nil, errors.New("get artist songs: no data")
	}

	a.patchSongInfo(ctx, artistSongs.Data.Info...)
	a.patchSongsInfo(ctx, artistSongs.Data.Info...)
	a.patchSongsLyric(ctx, artistSongs.Data.Info...)
	songs := translate(artistSongs.Data.Info...)
	return &api.Collection{
		Id:     strconv.Itoa(artistInfo.Data.SingerId),
		Name:   strings.TrimSpace(artistInfo.Data.SingerName),
		PicURL: strings.ReplaceAll(artistInfo.Data.ImgURL, "{size}", "480"),
		Songs:  songs,
	}, nil
}

// 获取歌手信息
func (a *API) GetArtistInfoRaw(ctx context.Context, singerId string) (*ArtistInfoResponse, error) {
	params := ghttp.Params{
		"singerid": singerId,
	}

	resp := new(ArtistInfoResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiGetArtistInfo)
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
		return nil, fmt.Errorf("get artist info: %s", resp.errorMessage())
	}

	return resp, nil
}

// 获取歌手歌曲，page: 页码；pageSize: 每页数量，-1获取全部
func (a *API) GetArtistSongsRaw(ctx context.Context, singerId string, page int, pageSize int) (*ArtistSongsResponse, error) {
	params := ghttp.Params{
		"singerid": singerId,
		"page":     page,
		"pagesize": pageSize,
	}

	resp := new(ArtistSongsResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiGetArtistSongs)
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
		return nil, fmt.Errorf("get artist songs: %s", resp.errorMessage())
	}

	return resp, nil
}
