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

func (a *API) GetArtist(ctx context.Context, artistId string) (*api.Collection, error) {
	_artistId, err := strconv.Atoi(artistId)
	if err != nil {
		return nil, err
	}

	resp, err := a.GetArtistRaw(ctx, _artistId)
	if err != nil {
		return nil, err
	}

	n := len(resp.HotSongs)
	if n == 0 {
		return nil, errors.New("get artist: no data")
	}

	a.patchSongsURL(ctx, songDefaultBR, resp.HotSongs...)
	a.patchSongsLyric(ctx, resp.HotSongs...)
	songs := translate(resp.HotSongs...)
	return &api.Collection{
		Id:     strconv.Itoa(resp.Artist.Id),
		Name:   strings.TrimSpace(resp.Artist.Name),
		PicURL: resp.Artist.PicURL,
		Songs:  songs,
	}, nil
}

// 获取歌手
func (a *API) GetArtistRaw(ctx context.Context, artistId int) (*ArtistResponse, error) {
	resp := new(ArtistResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodPost, fmt.Sprintf(apiGetArtist, artistId))
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
		return nil, fmt.Errorf("get artist: %s", resp.errorMessage())
	}

	return resp, nil
}
