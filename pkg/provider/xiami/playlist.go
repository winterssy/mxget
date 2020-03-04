package xiami

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
	"github.com/winterssy/mxget/pkg/utils"
)

func (a *API) GetPlaylist(ctx context.Context, playlistId string) (*api.Collection, error) {
	resp, err := a.GetPlaylistDetailRaw(ctx, playlistId, 1, songRequestLimit)
	if err != nil {
		return nil, err
	}

	n, _ := strconv.Atoi(resp.Data.Data.CollectDetail.SongCount)
	if n == 0 {
		return nil, errors.New("get playlist: no data")
	}

	_songs := resp.Data.Data.CollectDetail.Songs
	if n > songRequestLimit {
		allSongs := resp.Data.Data.CollectDetail.AllSongs
		queue := make(chan []*Song)
		wg := new(sync.WaitGroup)
		for i := songRequestLimit; i < n; i += songRequestLimit {
			if ctx.Err() != nil {
				break
			}

			songIds := allSongs[i:utils.Min(i+songRequestLimit, n)]
			wg.Add(1)
			go func() {
				resp, err := a.GetSongsRaw(ctx, songIds...)
				if err != nil {
					wg.Done()
					return
				}
				queue <- resp.Data.Data.Songs
			}()
		}

		go func() {
			for s := range queue {
				_songs = append(_songs, s...)
				wg.Done()
			}
		}()
		wg.Wait()
	}

	a.patchSongsLyric(ctx, _songs...)
	songs := translate(_songs...)
	return &api.Collection{
		Id:     resp.Data.Data.CollectDetail.ListId,
		Name:   strings.TrimSpace(resp.Data.Data.CollectDetail.CollectName),
		PicURL: resp.Data.Data.CollectDetail.CollectLogo,
		Songs:  songs,
	}, nil
}

// 获取歌单详情，包含歌单信息跟歌曲
func (a *API) GetPlaylistDetailRaw(ctx context.Context, playlistId string, page int, pageSize int) (*PlaylistDetailResponse, error) {
	token, err := a.getToken(apiGetPlaylistDetail)
	if err != nil {
		return nil, err
	}

	model := map[string]interface{}{
		"listId": playlistId,
		"pagingVO": map[string]int{
			"page":     page,
			"pageSize": pageSize,
		},
	}
	params := signPayload(token, model)

	resp := new(PlaylistDetailResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiGetPlaylistDetail)
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
		return nil, fmt.Errorf("get playlist detail: %s", err.Error())
	}

	return resp, nil
}

// 获取歌单歌曲，不包含歌单信息
func (a *API) GetPlaylistSongsRaw(ctx context.Context, playlistId string, page int, pageSize int) (*PlaylistSongsResponse, error) {
	token, err := a.getToken(apiGetPlaylistSongs)
	if err != nil {
		return nil, err
	}

	model := map[string]interface{}{
		"listId": playlistId,
		"pagingVO": map[string]int{
			"page":     page,
			"pageSize": pageSize,
		},
	}
	params := signPayload(token, model)

	resp := new(PlaylistSongsResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodGet, apiGetPlaylistSongs)
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
		return nil, fmt.Errorf("get playlist songs: %s", err.Error())
	}

	return resp, nil
}
