package migu

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/gjson"
	"github.com/winterssy/mxget/pkg/api"
)

func (a *API) SearchSongs(ctx context.Context, keyword string) ([]*api.Song, error) {
	resp, err := a.SearchSongsRawV2(ctx, keyword, 1, 50)
	if err != nil {
		return nil, err
	}

	musics := resp.GetArray("musics")
	n := len(musics)
	if n == 0 {
		return nil, errors.New("search songs: no data")
	}

	songs := make([]*api.Song, n)
	for i, s := range musics {
		obj := s.ToObject()
		songs[i] = &api.Song{
			Id:     obj.GetString("id"),
			Name:   obj.GetString("songName"),
			Artist: strings.ReplaceAll(obj.GetString("singerName"), ", ", "/"),
			Album:  obj.GetString("albumName"),
		}
	}
	return songs, nil
}

// 搜索歌曲
func (a *API) SearchSongsRawV1(ctx context.Context, keyword string, page int, pageSize int) (*SearchSongsResponse, error) {
	switchOption := map[string]int{
		"song":     1,
		"album":    0,
		"singer":   0,
		"tagSong":  0,
		"mvSong":   0,
		"songlist": 0,
		"bestShow": 0,
	}
	enc, _ := gjson.Marshal(switchOption)
	params := ghttp.Params{
		"searchSwitch": enc,
		"text":         keyword,
		"pageNo":       page,
		"pageSize":     pageSize,
	}

	resp := new(SearchSongsResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiSearchV1)
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
		return nil, fmt.Errorf("search songs: %s", resp.errorMessage())
	}

	return resp, nil
}

// 搜索歌曲
func (a *API) SearchSongsRawV2(ctx context.Context, keyword string, page int, pageSize int) (ghttp.H, error) {
	params := ghttp.Params{
		"keyword": keyword,
		"pgc":     page,
		"rows":    pageSize,
	}

	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiSearchV2)
	req.SetQuery(params)
	req.SetContext(ctx)
	r, err := a.SendRequest(req)
	if err != nil {
		return nil, err
	}

	data, err := r.H()
	if err != nil {
		return nil, err
	}

	if !data.GetBoolean("success") {
		return data, fmt.Errorf("search songs: %s", data.GetString("msg"))
	}

	return data, nil
}
