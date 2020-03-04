package kuwo

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
)

func (a *API) SearchSongs(ctx context.Context, keyword string) ([]*api.Song, error) {
	resp, err := a.SearchSongsRaw(ctx, keyword, 1, 50)
	if err != nil {
		return nil, err
	}

	n := len(resp.Data.List)
	if n == 0 {
		return nil, errors.New("search songs: no data")
	}

	songs := make([]*api.Song, n)
	for i, s := range resp.Data.List {
		songs[i] = &api.Song{
			Id:     strconv.Itoa(s.RId),
			Name:   strings.TrimSpace(s.Name),
			Artist: strings.TrimSpace(strings.ReplaceAll(s.Artist, "&", "/")),
			Album:  strings.TrimSpace(s.Album),
		}
	}
	return songs, nil
}

// 搜索歌曲
func (a *API) SearchSongsRaw(ctx context.Context, keyword string, page int, pageSize int) (*SearchSongsResponse, error) {
	params := ghttp.Params{
		"key": keyword,
		"pn":  page,
		"rn":  pageSize,
	}

	resp := new(SearchSongsResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiSearch)
	req.SetQuery(params)
	req.SetContext(ctx)
	r, err := a.SendRequest(req)
	if err == nil {
		err = r.JSON(resp)
	}
	if err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		if resp.Code == -1 {
			err = errors.New("search songs: no data")
		} else {
			err = fmt.Errorf("search songs: %s", resp.errorMessage())
		}
		return nil, err
	}

	return resp, nil
}
