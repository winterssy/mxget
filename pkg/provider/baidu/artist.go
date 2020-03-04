package baidu

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
)

func (a *API) GetArtist(ctx context.Context, tingUid string) (*api.Collection, error) {
	resp, err := a.GetArtistRaw(ctx, tingUid, 0, 50)
	if err != nil {
		return nil, err
	}

	n := len(resp.SongList)
	if n == 0 {
		return nil, errors.New("get artist: no data")
	}

	a.patchSongsURL(ctx, resp.SongList...)
	a.patchSongsLyric(ctx, resp.SongList...)
	songs := translate(resp.SongList...)
	return &api.Collection{
		Id:     resp.ArtistInfo.TingUid,
		Name:   strings.TrimSpace(resp.ArtistInfo.Name),
		PicURL: strings.SplitN(resp.ArtistInfo.AvatarBig, "@", 2)[0],
		Songs:  songs,
	}, nil
}

// 获取歌手
func (a *API) GetArtistRaw(ctx context.Context, tingUid string, offset int, limits int) (*ArtistResponse, error) {
	params := ghttp.Params{
		"tinguid": tingUid,
		"offset":  offset,
		"limits":  limits,
	}

	resp := new(ArtistResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiGetArtist)
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
		return nil, fmt.Errorf("get artist: %s", resp.errorMessage())
	}

	return resp, nil
}
