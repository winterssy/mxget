package kugou

import (
	"context"
	"strconv"
	"strings"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
	"github.com/winterssy/mxget/pkg/concurrency"
	"github.com/winterssy/mxget/pkg/request"
	"github.com/winterssy/mxget/pkg/utils"
)

const (
	apiSearch           = "http://mobilecdn.kugou.com/api/v3/search/song"
	apiGetSong          = "http://m.kugou.com/app/i/getSongInfo.php?cmd=playInfo"
	apiGetSongURL       = "http://trackercdn.kugou.com/i/v2/?pid=2&behavior=play&cmd=25"
	apiGetSongLyric     = "http://m.kugou.com/app/i/krc.php?cmd=100&timelength=1"
	apiGetArtistInfo    = "http://mobilecdn.kugou.com/api/v3/singer/info"
	apiGetArtistSongs   = "http://mobilecdn.kugou.com/api/v3/singer/song"
	apiGetAlbumInfo     = "http://mobilecdn.kugou.com/api/v3/album/info"
	apiGetAlbumSongs    = "http://mobilecdn.kugou.com/api/v3/album/song"
	apiGetPlaylistInfo  = "http://mobilecdn.kugou.com/api/v3/special/info"
	apiGetPlaylistSongs = "http://mobilecdn.kugou.com/api/v3/special/song"
)

var (
	std = New(request.DefaultClient)

	defaultHeaders = ghttp.Headers{
		"Origin":  "https://www.kugou.com",
		"Referer": "https://www.kugou.com",
	}
)

type (
	CommonResponse struct {
		Status  int    `json:"status"`
		Error   string `json:"error,omitempty"`
		ErrCode int    `json:"errcode"`
	}

	SearchSongsResponse struct {
		CommonResponse
		Data struct {
			Total int `json:"total"`
			Info  []*struct {
				Hash       string `json:"hash"`
				HQHash     string `json:"320hash"`
				SQHash     string `json:"sqhash"`
				SongName   string `json:"songname"`
				SingerName string `json:"singername"`
				AlbumId    string `json:"album_id"`
				AlbumName  string `json:"album_name"`
			} `json:"info"`
		} `json:"data"`
	}

	Song struct {
		Hash         string `json:"hash"`
		SongName     string `json:"songName"`
		SingerId     int    `json:"singerId"`
		SingerName   string `json:"singerName"`
		ChoricSinger string `json:"choricSinger"`
		FileName     string `json:"fileName"`
		ExtName      string `json:"extName"`
		AlbumId      int    `json:"albumid"`
		AlbumImg     string `json:"album_img"`
		Extra        struct {
			PQHash string `json:"128hash"`
			HQHash string `json:"320hash"`
			SQHash string `json:"sqhash"`
		} `json:"extra"`
		URL       string `json:"url"`
		AlbumName string `json:"-"`
		Lyric     string `json:"-"`
	}

	SongResponse struct {
		CommonResponse
		Song
	}

	SongURLResponse struct {
		Status  int      `json:"status"`
		BitRate int      `json:"bitRate"`
		ExtName string   `json:"extName"`
		URL     []string `json:"url"`
	}

	ArtistInfo struct {
		SingerId   int    `json:"singerid"`
		SingerName string `json:"singername"`
		ImgURL     string `json:"imgurl"`
	}

	ArtistInfoResponse struct {
		CommonResponse
		Data ArtistInfo `json:"data"`
	}

	ArtistSongsResponse struct {
		CommonResponse
		Data struct {
			Info []*Song `json:"info"`
		} `json:"data"`
	}

	AlbumInfo struct {
		AlbumId   int    `json:"albumid"`
		AlbumName string `json:"albumname"`
		ImgURL    string `json:"imgurl"`
	}

	AlbumInfoResponse struct {
		CommonResponse
		Data AlbumInfo `json:"data"`
	}

	AlbumSongsResponse struct {
		CommonResponse
		Data struct {
			Info []*Song `json:"info"`
		} `json:"data"`
	}

	PlaylistInfo struct {
		SpecialId   int    `json:"specialid"`
		SpecialName string `json:"specialname"`
		ImgURL      string `json:"imgurl"`
	}

	PlaylistInfoResponse struct {
		CommonResponse
		Data PlaylistInfo `json:"data"`
	}

	PlaylistSongsResponse struct {
		CommonResponse
		Data struct {
			Info []*Song `json:"info"`
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
	if c.Error == "" {
		return strconv.Itoa(c.ErrCode)
	}
	return c.Error
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

func (a *ArtistInfoResponse) String() string {
	return utils.PrettyJSON(a)
}

func (a *ArtistSongsResponse) String() string {
	return utils.PrettyJSON(a)
}

func (a *AlbumInfoResponse) String() string {
	return utils.PrettyJSON(a)
}

func (a *AlbumSongsResponse) String() string {
	return utils.PrettyJSON(a)
}

func (p *PlaylistInfoResponse) String() string {
	return utils.PrettyJSON(p)
}

func (p *PlaylistSongsResponse) String() string {
	return utils.PrettyJSON(p)
}

func (a *API) SendRequest(req *ghttp.Request) (*ghttp.Response, error) {
	req.SetHeaders(defaultHeaders)
	return a.Client.Do(req)
}

func (a *API) patchSongInfo(ctx context.Context, songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		if ctx.Err() != nil {
			break
		}

		c.Add(1)
		go func(s *Song) {
			resp, err := a.GetSongRaw(ctx, s.Hash)
			if err == nil {
				s.SongName = resp.SongName
				s.SingerId = resp.SingerId
				s.SingerName = resp.SingerName
				s.ChoricSinger = resp.ChoricSinger
				s.AlbumId = resp.AlbumId
				s.AlbumImg = resp.AlbumImg
				s.Extra = resp.Extra
				s.URL = resp.URL
			}
			c.Done()
		}(s)
	}
	c.Wait()
}

func (a *API) patchSongsInfo(ctx context.Context, songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		if ctx.Err() != nil {
			break
		}

		c.Add(1)
		go func(s *Song) {
			if s.AlbumId != 0 {
				resp, err := a.GetAlbumInfoRaw(ctx, strconv.Itoa(s.AlbumId))
				if err == nil {
					s.AlbumName = resp.Data.AlbumName
				}
			}
			c.Done()
		}(s)
	}
	c.Wait()
}

func (a *API) patchSongsURL(ctx context.Context, songs ...*Song) {
	c := concurrency.New(32)
	for _, s := range songs {
		if ctx.Err() != nil {
			break
		}

		if s.URL != "" {
			continue
		}
		c.Add(1)
		go func(s *Song) {
			url, err := a.GetSongURL(ctx, s.Hash)
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
			lyric, err := a.GetSongLyric(ctx, s.Hash)
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
			Id:        s.Hash,
			Name:      strings.TrimSpace(s.SongName),
			Artist:    strings.TrimSpace(strings.ReplaceAll(s.ChoricSinger, "、", "/")),
			Album:     strings.TrimSpace(s.AlbumName),
			PicURL:    strings.ReplaceAll(s.AlbumImg, "{size}", "480"),
			Lyric:     s.Lyric,
			ListenURL: s.URL,
		}
	}
	return songs
}
