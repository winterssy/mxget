package netease_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winterssy/mxget/pkg/provider/netease"
)

var (
	client *netease.API
	ctx    context.Context
)

func setup() {
	client = netease.Client()
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
	song, err := client.GetSong(ctx, "444269135")
	if assert.NoError(t, err) {
		t.Log(song)
	}
}

func TestAPI_GetSongURL(t *testing.T) {
	url, err := client.GetSongURL(ctx, 444269135, 320)
	if assert.NoError(t, err) {
		t.Log(url)
	}
}

func TestAPI_GetSongLyric(t *testing.T) {
	lyric, err := client.GetSongLyric(ctx, 444269135)
	if assert.NoError(t, err) {
		t.Log(lyric)
	}
}

func TestAPI_GetArtist(t *testing.T) {
	artist, err := client.GetArtist(ctx, "1045123")
	if assert.NoError(t, err) {
		t.Log(artist)
	}
}

func TestAPI_GetAlbum(t *testing.T) {
	album, err := client.GetAlbum(ctx, "35023284")
	if assert.NoError(t, err) {
		t.Log(album)
	}
}

func TestAPI_GetPlaylist(t *testing.T) {
	playlist, err := client.GetPlaylist(ctx, "156934569")
	if assert.NoError(t, err) {
		t.Log(playlist)
	}
}
