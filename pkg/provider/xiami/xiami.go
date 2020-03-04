package xiami

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
	"github.com/winterssy/mxget/pkg/concurrency"
	"github.com/winterssy/mxget/pkg/request"
	"github.com/winterssy/mxget/pkg/utils"
)

const (
	apiSearch            = "https://acs.m.xiami.com/h5/mtop.alimusic.search.searchservice.searchsongs/1.0/?appKey=23649156"
	apiGetSongDetail     = "https://acs.m.xiami.com/h5/mtop.alimusic.music.songservice.getsongdetail/1.0/?appKey=23649156"
	apiGetSongLyric      = "https://acs.m.xiami.com/h5/mtop.alimusic.music.lyricservice.getsonglyrics/1.0/?appKey=23649156"
	apiGetSongs          = "https://acs.m.xiami.com/h5/mtop.alimusic.music.songservice.getsongs/1.0/?appKey=23649156"
	apiGetArtistInfo     = "https://acs.m.xiami.com/h5/mtop.alimusic.music.artistservice.getartistdetail/1.0/?appKey=23649156"
	apiGetArtistSongs    = "https://acs.m.xiami.com/h5/mtop.alimusic.music.songservice.getartistsongs/1.0/?appKey=23649156"
	apiGetAlbum          = "https://acs.m.xiami.com/h5/mtop.alimusic.music.albumservice.getalbumdetail/1.0/?appKey=23649156"
	apiGetPlaylistDetail = "https://h5api.m.xiami.com/h5/mtop.alimusic.music.list.collectservice.getcollectdetail/1.0/?appKey=23649156"
	apiGetPlaylistSongs  = "https://h5api.m.xiami.com/h5/mtop.alimusic.music.list.collectservice.getcollectsongs/1.0/?appKey=23649156"
	apiLogin             = "https://h5api.m.xiami.com/h5/mtop.alimusic.xuser.facade.xiamiuserservice.login/1.0/?appKey=23649156"

	songRequestLimit = 200
)

var (
	std = New(request.DefaultClient)

	defaultHeaders = ghttp.Headers{
		"Origin":  "https://h.xiami.com",
		"Referer": "https://h.xiami.com",
	}
)

type (
	CommonResponse struct {
		API string   `json:"api"`
		Ret []string `json:"ret"`
	}

	PagingVO struct {
		Count    string `json:"count"`
		Page     string `json:"page"`
		PageSize string `json:"pageSize"`
		Pages    string `json:"pages"`
	}

	ListenFile struct {
		Expire     string `json:"expire,omitempty"`
		FileSize   string `json:"fileSize"`
		Format     string `json:"format"`
		ListenFile string `json:"listenFile,omitempty"`
		Quality    string `json:"quality"`
		URL        string `json:"url,omitempty"`
	}

	Song struct {
		Album
		SongId       string       `json:"songId"`
		SongStringId string       `json:"songStringId"`
		SongName     string       `json:"songName"`
		Singers      string       `json:"singers"`
		SingerVOs    []Artist     `json:"singerVOs"`
		ListenFile   string       `json:"listenFile,omitempty"`
		ListenFiles  []ListenFile `json:"listenFiles"`
		Lyric        string       `json:"-"`
	}

	SearchSongsResponse struct {
		CommonResponse
		Data struct {
			Data struct {
				PagingVO PagingVO `json:"pagingVO"`
				Songs    []*Song  `json:"songs"`
			} `json:"data"`
		} `json:"data"`
	}

	SongDetailResponse struct {
		CommonResponse
		Data struct {
			Data struct {
				SongDetail Song `json:"songDetail"`
			} `json:"data"`
		} `json:"data"`
	}

	SongLyricResponse struct {
		CommonResponse
		Data struct {
			Data struct {
				Lyrics []struct {
					Content      string `json:"content"`
					FlagOfficial string `json:"flagOfficial"`
					LyricURL     string `json:"lyricUrl"`
					Type         string `json:"type"`
				} `json:"lyrics"`
			} `json:"data"`
		} `json:"data"`
	}

	SongsResponse struct {
		CommonResponse
		Data struct {
			Data struct {
				Songs []*Song `json:"songs"`
			} `json:"data"`
		} `json:"data"`
	}

	Artist struct {
		ArtistId       string `json:"artistId"`
		ArtistStringId string `json:"artistStringId"`
		ArtistName     string `json:"artistName"`
		ArtistLogo     string `json:"artistLogo"`
	}

	ArtistInfoResponse struct {
		CommonResponse
		Data struct {
			Data struct {
				ArtistDetailVO Artist `json:"artistDetailVO"`
			} `json:"data"`
		} `json:"data"`
	}

	ArtistSongsResponse struct {
		CommonResponse
		Data struct {
			Data struct {
				PagingVO PagingVO `json:"pagingVO"`
				Songs    []*Song  `json:"songs"`
			} `json:"data"`
		} `json:"data"`
	}

	Album struct {
		AlbumId       string `json:"albumId"`
		AlbumStringId string `json:"albumStringId"`
		AlbumName     string `json:"albumName"`
		AlbumLogo     string `json:"albumLogo"`
	}

	AlbumResponse struct {
		CommonResponse
		Data struct {
			Data struct {
				AlbumDetail struct {
					Album
					Songs []*Song `json:"songs"`
				} `json:"albumDetail"`
			} `json:"data"`
		} `json:"data"`
	}

	PlaylistDetailResponse struct {
		CommonResponse
		Data struct {
			Data struct {
				CollectDetail struct {
					ListId      string   `json:"listId"`
					CollectName string   `json:"collectName"`
					CollectLogo string   `json:"collectLogo"`
					SongCount   string   `json:"songCount"`
					AllSongs    []string `json:"allSongs"`
					Songs       []*Song  `json:"songs"`
					PagingVO    PagingVO `json:"pagingVO"`
				} `json:"collectDetail"`
			} `json:"data"`
		} `json:"data"`
	}

	PlaylistSongsResponse struct {
		CommonResponse
		Data struct {
			Data struct {
				Songs    []*Song  `json:"songs"`
				PagingVO PagingVO `json:"pagingVO"`
			} `json:"data"`
		} `json:"data"`
	}

	LoginResponse struct {
		CommonResponse
		Data struct {
			Data struct {
				AccessToken    string `json:"accessToken"`
				Expires        string `json:"expires"`
				NickName       string `json:"nickName"`
				RefreshExpires string `json:"refreshExpires"`
				RefreshToken   string `json:"refreshToken"`
				UserId         string `json:"userId"`
			} `json:"data"`
		} `json:"data"`
	}

	API struct {
		Client *ghttp.Client
	}
)

func New(client *ghttp.Client) *API {
	return &API{
		Client: client,
	}
}

func Client() *API {
	return std
}

func (c *CommonResponse) check() error {
	for _, s := range c.Ret {
		if strings.HasPrefix(s, "FAIL") {
			return errors.New(s)
		}
	}
	return nil
}

func (s *SearchSongsResponse) String() string {
	return utils.PrettyJSON(s)
}

func (s *SongDetailResponse) String() string {
	return utils.PrettyJSON(s)
}

func (s *SongsResponse) String() string {
	return utils.PrettyJSON(s)
}

func (a *ArtistInfoResponse) String() string {
	return utils.PrettyJSON(a)
}

func (a *ArtistSongsResponse) String() string {
	return utils.PrettyJSON(a)
}

func (a *AlbumResponse) String() string {
	return utils.PrettyJSON(a)
}

func (p *PlaylistDetailResponse) String() string {
	return utils.PrettyJSON(p)
}

func (p *PlaylistSongsResponse) String() string {
	return utils.PrettyJSON(p)
}

func (l *LoginResponse) String() string {
	return utils.PrettyJSON(l)
}

func (a *API) SendRequest(req *ghttp.Request) (*ghttp.Response, error) {
	req.SetHeaders(defaultHeaders)
	return a.Client.Do(req)
}

func (a *API) getToken(url string) (string, error) {
	const XiaMiToken = "_m_h5_tk"
	token, err := a.Client.Cookie(url, XiaMiToken)
	if err != nil {
		// 如果在cookie jar中没有找到对应cookie，发送预请求获取
		req, _ := ghttp.NewRequest(ghttp.MethodGet, url)
		var resp *ghttp.Response
		resp, err = a.SendRequest(req)
		if err == nil {
			token, err = resp.Cookie(XiaMiToken)
		}
	}

	if err != nil {
		return "", fmt.Errorf("can't get token: %s", err.Error())
	}

	return strings.Split(token.Value, "_")[0], nil
}

func (a *API) patchSongsLyric(ctx context.Context, songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		if ctx.Err() != nil {
			break
		}

		c.Add(1)
		go func(s *Song) {
			lyric, err := a.GetSongLyric(ctx, s.SongId)
			if err == nil {
				s.Lyric = lyric
			}
			c.Done()
		}(s)
	}
	c.Wait()
}

func songURL(listenFiles []ListenFile) string {
	for _, i := range listenFiles {
		if i.Quality == "l" {
			return i.URL + i.ListenFile
		}
	}
	return ""
}

func translate(src ...*Song) []*api.Song {
	songs := make([]*api.Song, len(src))
	for i, s := range src {
		url := songURL(s.ListenFiles)
		songs[i] = &api.Song{
			Id:        s.SongId,
			Name:      strings.TrimSpace(s.SongName),
			Artist:    strings.TrimSpace(strings.ReplaceAll(s.Singers, " / ", "/")),
			Album:     strings.TrimSpace(s.AlbumName),
			PicURL:    s.AlbumLogo,
			Lyric:     s.Lyric,
			ListenURL: url,
		}
	}
	return songs
}
