package baidu

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

	n := len(resp.Result.SongInfo.SongList)
	if n == 0 {
		return nil, errors.New("search songs: no data")
	}

	songs := make([]*api.Song, n)
	for i, s := range resp.Result.SongInfo.SongList {
		songs[i] = &api.Song{
			Id:     s.SongId,
			Name:   strings.TrimSpace(s.Title),
			Artist: strings.TrimSpace(strings.ReplaceAll(s.Author, ",", "/")),
			Album:  strings.TrimSpace(s.AlbumTitle),
		}
	}
	return songs, nil
}

// 搜索歌曲
func (a *API) SearchSongsRaw(ctx context.Context, keyword string, page int, pageSize int) (*SearchSongsResponse, error) {
	params := ghttp.Params{
		"query":     keyword,
		"page_no":   strconv.Itoa(page),
		"page_size": strconv.Itoa(pageSize),
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

	if resp.ErrorCode != 22000 {
		return nil, fmt.Errorf("search songs: %s", resp.errorMessage())
	}

	return resp, nil
}
