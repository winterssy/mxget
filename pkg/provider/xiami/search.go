package xiami

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
)

func (a *API) SearchSongs(ctx context.Context, keyword string) ([]*api.Song, error) {
	resp, err := a.SearchSongsRaw(ctx, keyword, 1, 50)
	if err != nil {
		return nil, err
	}

	n := len(resp.Data.Data.Songs)
	if n == 0 {
		return nil, errors.New("search songs: no data")
	}

	songs := make([]*api.Song, n)
	for i, s := range resp.Data.Data.Songs {
		songs[i] = &api.Song{
			Id:     s.SongId,
			Name:   strings.TrimSpace(s.SongName),
			Artist: strings.TrimSpace(strings.ReplaceAll(s.Singers, " / ", "/")),
			Album:  strings.TrimSpace(s.AlbumName),
		}
	}
	return songs, nil
}

// 搜索歌曲
func (a *API) SearchSongsRaw(ctx context.Context, keyword string, page int, pageSize int) (*SearchSongsResponse, error) {
	token, err := a.getToken(apiSearch)
	if err != nil {
		return nil, err
	}

	model := map[string]interface{}{
		"key": keyword,
		"pagingVO": map[string]int{
			"page":     page,
			"pageSize": pageSize,
		},
	}
	params := signPayload(token, model)

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

	err = resp.check()
	if err != nil {
		return nil, fmt.Errorf("search songs: %s", err.Error())
	}

	return resp, nil
}
