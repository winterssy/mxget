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

func (a *API) GetArtist(ctx context.Context, artistId string) (*api.Collection, error) {
	artistInfo, err := a.GetArtistInfoRaw(ctx, artistId)
	if err != nil {
		return nil, err
	}

	artistSongs, err := a.GetArtistSongsRaw(ctx, artistId, 1, 50)
	if err != nil {
		return nil, err
	}

	_songs := artistSongs.Data.Data.Songs
	n := len(_songs)
	if n == 0 {
		return nil, errors.New("get artist songs: no data")
	}

	a.patchSongsLyric(ctx, _songs...)
	songs := translate(_songs...)
	return &api.Collection{
		Id:     artistInfo.Data.Data.ArtistDetailVO.ArtistId,
		Name:   strings.TrimSpace(artistInfo.Data.Data.ArtistDetailVO.ArtistName),
		PicURL: artistInfo.Data.Data.ArtistDetailVO.ArtistLogo,
		Songs:  songs,
	}, nil
}

// 获取歌手信息
func (a *API) GetArtistInfoRaw(ctx context.Context, artistId string) (*ArtistInfoResponse, error) {
	token, err := a.getToken(apiGetArtistInfo)
	if err != nil {
		return nil, err
	}

	model := make(map[string]string)
	_, err = strconv.Atoi(artistId)
	if err != nil {
		model["artistStringId"] = artistId
	} else {
		model["artistId"] = artistId
	}
	params := signPayload(token, model)

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

	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("get artist info: %s", err.Error())
	}

	return resp, nil
}

// 获取歌手歌曲
func (a *API) GetArtistSongsRaw(ctx context.Context, artistId string, page int, pageSize int) (*ArtistSongsResponse, error) {
	token, err := a.getToken(apiGetArtistSongs)
	if err != nil {
		return nil, err
	}

	model := map[string]interface{}{
		"pagingVO": map[string]int{
			"page":     page,
			"pageSize": pageSize,
		},
	}
	_, err = strconv.Atoi(artistId)
	if err != nil {
		model["artistStringId"] = artistId
	} else {
		model["artistId"] = artistId
	}
	params := signPayload(token, model)

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

	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("get artist songs: %s", err.Error())
	}

	return resp, nil
}
