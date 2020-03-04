package kuwo

import (
	"context"
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
	apiSearch         = "http://www.kuwo.cn/api/www/search/searchMusicBykeyWord"
	apiGetSong        = "http://www.kuwo.cn/api/www/music/musicInfo"
	apiGetSongURL     = "http://www.kuwo.cn/url?format=mp3&response=url&type=convert_url3"
	apiGetSongLyric   = "http://www.kuwo.cn/newh5/singles/songinfoandlrc"
	apiGetArtistInfo  = "http://www.kuwo.cn/api/www/artist/artist"
	apiGetArtistSongs = "http://www.kuwo.cn/api/www/artist/artistMusic"
	apiGetAlbum       = "http://www.kuwo.cn/api/www/album/albumInfo"
	apiGetPlaylist    = "http://www.kuwo.cn/api/www/playlist/playListInfo"

	songDefaultBR = 128
)

var (
	std = New(request.DefaultClient)
)

type (
	CommonResponse struct {
		Code int    `json:"code"`
		Msg  string `json:"msg,omitempty"`
	}

	Song struct {
		RId             int    `json:"rid"`
		Name            string `json:"name"`
		ArtistId        int    `json:"artistid"`
		Artist          string `json:"artist"`
		AlbumId         int    `json:"albumid"`
		Album           string `json:"album"`
		AlbumPic        string `json:"albumpic"`
		Track           int    `json:"track"`
		IsListenFee     bool   `json:"isListenFee"`
		SongTimeMinutes string `json:"songTimeMinutes"`
		Lyric           string `json:"-"`
		URL             string `json:"-"`
	}

	SearchSongsResponse struct {
		CommonResponse
		Data struct {
			Total string  `json:"total"`
			List  []*Song `json:"list"`
		} `json:"data"`
	}

	SongResponse struct {
		CommonResponse
		Data Song `json:"data"`
	}

	SongURLResponse struct {
		CommonResponse
		URL string `json:"url"`
	}

	SongLyricResponse struct {
		Status int    `json:"status"`
		Msg    string `json:"msg,omitempty"`
		Data   struct {
			LrcList []struct {
				Time      string `json:"time"`
				LineLyric string `json:"lineLyric"`
			} `json:"lrclist"`
		} `json:"data"`
	}

	ArtistInfo struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		Pic300 string `json:"pic300"`
	}

	ArtistInfoResponse struct {
		CommonResponse
		Data ArtistInfo `json:"data"`
	}

	ArtistSongsResponse struct {
		CommonResponse
		Data struct {
			List []*Song `json:"list"`
		} `json:"data"`
	}

	AlbumResponse struct {
		CommonResponse
		Data struct {
			AlbumId   int     `json:"albumId"`
			Album     string  `json:"album"`
			Pic       string  `json:"pic"`
			MusicList []*Song `json:"musicList"`
		} `json:"data"`
	}

	PlaylistResponse struct {
		CommonResponse
		Data struct {
			Id        int     `json:"id"`
			Name      string  `json:"name"`
			Img700    string  `json:"img700"`
			MusicList []*Song `json:"musicList"`
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

func (c *CommonResponse) errorMessage() string {
	if c.Msg == "" {
		return strconv.Itoa(c.Code)
	}
	return c.Msg
}

func (s *SearchSongsResponse) String() string {
	return utils.PrettyJSON(s)
}

func (s *SongResponse) String() string {
	return utils.PrettyJSON(s)
}

func (s *SongURLResponse) String() string {
	return utils.PrettyJSON(s)
}

func (s *SongLyricResponse) String() string {
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

func (p *PlaylistResponse) String() string {
	return utils.PrettyJSON(p)
}

func (a *API) SendRequest(req *ghttp.Request) (*ghttp.Response, error) {
	// csrf 必须跟 kw_token 保持一致
	csrf := "0"
	cookie, err := a.Client.Cookie(req.URL.String(), "kw_token")
	if err != nil {
		req.AddCookie(&http.Cookie{
			Name:  "kw_token",
			Value: csrf,
		})
	} else {
		csrf = cookie.Value
	}

	headers := ghttp.Headers{
		"Origin":  "http://www.kuwo.cn",
		"Referer": "http://www.kuwo.cn",
		"csrf":    csrf,
	}
	req.SetHeaders(headers)
	return a.Client.Do(req)
}

func (a *API) patchSongsURL(ctx context.Context, br int, songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		if ctx.Err() != nil {
			break
		}

		c.Add(1)
		go func(s *Song) {
			url, err := a.GetSongURL(ctx, s.RId, br)
			if err == nil {
				s.URL = url
			}
			c.Done()
		}(s)
	}
	c.Wait()
}

func (a *API) patchSongsLyric(ctx context.Context, songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		if ctx.Err() != nil {
			break
		}

		c.Add(1)
		go func(s *Song) {
			lyric, err := a.GetSongLyric(ctx, s.RId)
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
		songs[i] = &api.Song{
			Id:        strconv.Itoa(s.RId),
			Name:      strings.TrimSpace(s.Name),
			Artist:    strings.TrimSpace(strings.ReplaceAll(s.Artist, "&", "/")),
			Album:     strings.TrimSpace(s.Album),
			PicURL:    s.AlbumPic,
			Lyric:     s.Lyric,
			ListenURL: s.URL,
		}
	}
	return songs
}
