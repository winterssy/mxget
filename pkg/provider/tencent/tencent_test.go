package tencent_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winterssy/mxget/pkg/provider/tencent"
)

var (
	client *tencent.API
	ctx    context.Context
)

func setup() {
	client = tencent.Client()
	ctx = context.Background()
}

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func TestAPI_SearchSongs(t *testing.T) {
	result, err := client.SearchSongs(ctx, "Alan Walker")
	if assert.NoError(t, err) {
		t.Log(result)
	}
}

func TestAPI_GetSong(t *testing.T) {
	song, err := client.GetSong(ctx, "0015xqSa3G50mK")
	if assert.NoError(t, err) {
		t.Log(song)
	}
}

func TestAPI_GetSongURLV1(t *testing.T) {
	url, err := client.GetSongURLV1(ctx, "0015xqSa3G50mK", "0015xqSa3G50mK")
	if assert.NoError(t, err) {
		t.Log(url)
	}
}

func TestAPI_GetSongURLV2(t *testing.T) {
	url, err := client.GetSongURLV2(ctx, "0015xqSa3G50mK")
	if assert.NoError(t, err) {
		t.Log(url)
	}
}

func TestAPI_GetSongLyric(t *testing.T) {
	lyric, err := client.GetSongLyric(ctx, "002Zkt5S2z8JZx")
	if assert.NoError(t, err) {
		t.Log(lyric)
	}
}

func TestAPI_GetArtist(t *testing.T) {
	artist, err := client.GetArtist(ctx, "000Sp0Bz4JXH0o")
	if assert.NoError(t, err) {
		t.Log(artist)
	}
}

func TestAPI_GetAlbum(t *testing.T) {
	album, err := client.GetAlbum(ctx, "002fRO0N4FftzY")
	if assert.NoError(t, err) {
		t.Log(album)
	}
}

func TestAPI_GetPlaylist(t *testing.T) {
	playlist, err := client.GetPlaylist(ctx, "5474239760")
	if assert.NoError(t, err) {
		t.Log(playlist)
	}
}
