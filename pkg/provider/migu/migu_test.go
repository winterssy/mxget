package migu_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winterssy/mxget/pkg/provider/migu"
)

var (
	client *migu.API
	ctx    context.Context
)

func setup() {
	client = migu.Client()
	ctx = context.Background()
}

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func TestAPI_SearchSongs(t *testing.T) {
	result, err := client.SearchSongs(ctx, "周杰伦")
	if assert.NoError(t, err) {
		t.Log(result)
	}
}

func TestAPI_GetSong(t *testing.T) {
	song, err := client.GetSong(ctx, "63273402938")
	if assert.NoError(t, err) {
		t.Log(song)
	}
}

func TestAPI_GetSongURL(t *testing.T) {
	resp, err := client.GetSongURL(ctx, "600908000002677565", "2")
	if assert.NoError(t, err) {
		t.Log(resp)
	}
}

func TestAPI_GetSongLyric(t *testing.T) {
	lyric, err := client.GetSongLyric(ctx, "63273402938")
	if assert.NoError(t, err) {
		t.Log(lyric)
	}
}

func TestAPI_GetSongPic(t *testing.T) {
	pic, err := client.GetSongPic(ctx, "1121439251")
	if assert.NoError(t, err) {
		t.Log(pic)
	}
}

func TestAPI_GetArtist(t *testing.T) {
	artist, err := client.GetArtist(ctx, "112")
	if assert.NoError(t, err) {
		t.Log(artist)
	}
}

func TestAPI_GetAlbum(t *testing.T) {
	album, err := client.GetAlbum(ctx, "1121438701")
	if assert.NoError(t, err) {
		t.Log(album)
	}
}

func TestAPI_GetPlaylist(t *testing.T) {
	playlist, err := client.GetPlaylist(ctx, "159248239")
	if assert.NoError(t, err) {
		t.Log(playlist)
	}
}
