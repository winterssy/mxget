package netease

import (
	"context"
	"encoding/hex"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
	"github.com/winterssy/mxget/pkg/concurrency"
	"github.com/winterssy/mxget/pkg/request"
	"github.com/winterssy/mxget/pkg/utils"
)

const (
	linuxAPI          = "https://music.163.com/api/linux/forward"
	apiSearch         = "https://music.163.com/weapi/search/get"
	apiGetSongs       = "https://music.163.com/weapi/v3/song/detail"
	apiGetSongsURL    = "https://music.163.com/weapi/song/enhance/player/url"
	apiGetArtist      = "https://music.163.com/weapi/v1/artist/%d"
	apiGetAlbum       = "https://music.163.com/weapi/v1/album/%d"
	apiGetPlaylist    = "https://music.163.com/weapi/v3/playlist/detail"
	apiEmailLogin     = "https://music.163.com/weapi/login"
	apiCellphoneLogin = "https://music.163.com/weapi/login/cellphone"
	apiRefreshLogin   = "https://music.163.com/weapi/login/token/refresh"
	apiLogout         = "https://music.163.com/weapi/logout"

	songRequestLimit = 1000
	songDefaultBR    = 128
)

var (
	std = New(request.DefaultClient)

	defaultCookie = createCookie()

	defaultHeaders = ghttp.Headers{
		"Origin":  "https://music.163.com",
		"Referer": "https://music.163.com",
	}

	defaultCookieList = map[string]string{
		"os":            "linux",
		"appver":        "1.0.0.1026",
		"osver":         "Ubuntu%2016.10",
		"MUSIC_U":       "78d411095f4b022667bc8ec49e9a44cca088df057d987f5feaf066d37458e41c4a7d9447977352cf27ea9fee03f6ec4441049cea1c6bb9b6",
		"__remember_me": "true",
	}
)

type (
	CommonResponse struct {
		Code int    `json:"code"`
		Msg  string `json:"msg,omitempty"`
	}

	Song struct {
		Id      int      `json:"id"`
		Name    string   `json:"name"`
		Artists []Artist `json:"ar"`
		Album   Album    `json:"al"`
		Track   int      `json:"no"`
		Lyric   string   `json:"-"`
		URL     string   `json:"-"`
	}

	SearchSongsResponse struct {
		CommonResponse
		Result struct {
			Songs []*struct {
				Id      int    `json:"id"`
				Name    string `json:"name"`
				Artists []struct {
					Id   int    `json:"id"`
					Name string `json:"name"`
				} `json:"artists"`
				Album struct {
					Id   int    `json:"id"`
					Name string `json:"name"`
				} `json:"album"`
			} `json:"songs"`
			SongCount int `json:"songCount"`
		} `json:"result"`
	}

	SongURL struct {
		Id   int    `json:"id"`
		URL  string `json:"url"`
		BR   int    `json:"br"`
		Code int    `json:"code"`
	}

	SongsResponse struct {
		CommonResponse
		Songs []*Song `json:"songs"`
	}

	SongURLResponse struct {
		CommonResponse
		Data []struct {
			Code int    `json:"code"`
			Id   int    `json:"id"`
			BR   int    `json:"br"`
			URL  string `json:"url"`
		} `json:"data"`
	}

	SongLyricResponse struct {
		CommonResponse
		Lrc struct {
			Lyric string `json:"lyric"`
		} `json:"lrc"`
		TLyric struct {
			Lyric string `json:"lyric"`
		} `json:"tlyric"`
	}

	Artist struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		PicURL string `json:"picUrl"`
	}

	ArtistResponse struct {
		CommonResponse
		Artist struct {
			Artist
		} `json:"artist"`
		HotSongs []*Song `json:"hotSongs"`
	}

	Album struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		PicURL string `json:"picUrl"`
	}

	AlbumResponse struct {
		CommonResponse
		Album Album   `json:"album"`
		Songs []*Song `json:"songs"`
	}

	PlaylistResponse struct {
		CommonResponse
		Playlist struct {
			Id       int     `json:"id"`
			Name     string  `json:"name"`
			PicURL   string  `json:"coverImgUrl"`
			Tracks   []*Song `json:"tracks"`
			TrackIds []struct {
				Id int `json:"id"`
			} `json:"trackIds"`
			Total int `json:"trackCount"`
		} `json:"playlist"`
	}

	LoginResponse struct {
		CommonResponse
		LoginType int `json:"loginType"`
		Account   struct {
			Id       int    `json:"id"`
			UserName string `json:"userName"`
		} `json:"account"`
	}

	API struct {
		Client *ghttp.Client
	}
)

func createCookie() *http.Cookie {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return &http.Cookie{
		Name:  "_ntes_nuid",
		Value: hex.EncodeToString(b),
	}
}

func New(client *ghttp.Client) *API {
	return &API{
		Client: client,
	}
}

func Client() *API {
	return std
}

func (c *CommonResponse) errorMessage() string {
	if c.Msg == "" {
		return strconv.Itoa(c.Code)
	}
	return c.Msg
}

func (s *SearchSongsResponse) String() string {
	return utils.PrettyJSON(s)
}

func (s *SongsResponse) String() string {
	return utils.PrettyJSON(s)
}

func (s *SongURLResponse) String() string {
	return utils.PrettyJSON(s)
}

func (s *SongLyricResponse) String() string {
	return utils.PrettyJSON(s)
}

func (a *ArtistResponse) String() string {
	return utils.PrettyJSON(a)
}

func (a *AlbumResponse) String() string {
	return utils.PrettyJSON(a)
}

func (p *PlaylistResponse) String() string {
	return utils.PrettyJSON(p)
}

func (e *LoginResponse) String() string {
	return utils.PrettyJSON(e)
}

func (a *API) SendRequest(req *ghttp.Request) (*ghttp.Response, error) {
	req.SetHeaders(defaultHeaders)

	// 如果已经登录，不需要额外设置cookies，cookie jar会自动管理
	_, err := a.Client.Cookie(req.URL.String(), "MUSIC_U")
	if err != nil {
		req.AddCookie(defaultCookie)
	}

	for key, value := range defaultCookieList {
		req.AddCookie(&http.Cookie{
			Name:  key,
			Value: value,
		})
	}
	return a.Client.Do(req)
}

func (a *API) patchSongsURL(ctx context.Context, br int, songs ...*Song) {
	ids := make([]int, len(songs))
	for i, s := range songs {
		ids[i] = s.Id
	}

	resp, err := a.GetSongsURLRaw(ctx, br, ids...)
	if err == nil && len(resp.Data) != 0 {
		m := make(map[int]string, len(resp.Data))
		for _, i := range resp.Data {
			m[i.Id] = i.URL
		}
		for _, s := range songs {
			s.URL = m[s.Id]
		}
	}
}

func (a *API) patchSongsLyric(ctx context.Context, songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		if ctx.Err() != nil {
			break
		}

		c.Add(1)
		go func(s *Song) {
			lyric, err := a.GetSongLyric(ctx, s.Id)
			if err == nil {
				s.Lyric = lyric
			}
			c.Done()
		}(s)
	}
	c.Wait()
}

func translate(src ...*Song) []*api.Song {
	songs := make([]*api.Song, len(src))
	for i, s := range src {
		artists := make([]string, len(s.Artists))
		for j, a := range s.Artists {
			artists[j] = strings.TrimSpace(a.Name)
		}
		songs[i] = &api.Song{
			Id:        strconv.Itoa(s.Id),
			Name:      strings.TrimSpace(s.Name),
			Artist:    strings.Join(artists, "/"),
			Album:     strings.TrimSpace(s.Album.Name),
			PicURL:    s.Album.PicURL,
			Lyric:     s.Lyric,
			ListenURL: s.URL,
		}
	}
	return songs
}
