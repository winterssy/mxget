package kugou_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winterssy/mxget/pkg/provider/kugou"
)

var (
	client *kugou.API
	ctx    context.Context
)

func setup() {
	client = kugou.Client()
	ctx = context.Background()
}

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func TestAPI_SearchSongs(t *testing.T) {
	result, err := client.SearchSongs(ctx, "五月天")
	if assert.NoError(t, err) {
		t.Log(result)
	}
}

func TestAPI_GetSong(t *testing.T) {
	song, err := client.GetSong(ctx, "1571941D82D63AD614E35EAD9DB6A6A2")
	if assert.NoError(t, err) {
		t.Log(song)
	}
}

func TestAPI_GetSongURL(t *testing.T) {
	url, err := client.GetSongURL(ctx, "1571941D82D63AD614E35EAD9DB6A6A2")
	if assert.NoError(t, err) {
		t.Log(url)
	}
}

func TestAPI_GetSongLyric(t *testing.T) {
	lyric, err := client.GetSongLyric(ctx, "1571941D82D63AD614E35EAD9DB6A6A2")
	if assert.NoError(t, err) {
		t.Log(lyric)
	}
}

func TestAPI_GetArtist(t *testing.T) {
	artist, err := client.GetArtist(ctx, "8965")
	if assert.NoError(t, err) {
		t.Log(artist)
	}
}

func TestAPI_GetAlbum(t *testing.T) {
	album, err := client.GetAlbum(ctx, "976965")
	if assert.NoError(t, err) {
		t.Log(album)
	}
}

func TestAPI_GetPlaylist(t *testing.T) {
	playlist, err := client.GetPlaylist(ctx, "610433")
	if assert.NoError(t, err) {
		t.Log(playlist)
	}
}
