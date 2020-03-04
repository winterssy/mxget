package cmd

import (
	"github.com/spf13/cobra"
	"github.com/winterssy/mxget/cmd/album"
	"github.com/winterssy/mxget/cmd/artist"
	"github.com/winterssy/mxget/cmd/config"
	"github.com/winterssy/mxget/cmd/playlist"
	"github.com/winterssy/mxget/cmd/search"
	"github.com/winterssy/mxget/cmd/serve"
	"github.com/winterssy/mxget/cmd/song"
	"github.com/winterssy/mxget/internal/settings"
)

const (
	Version = "1.0.0"
)

var CmdRoot = &cobra.Command{
	Use:   "mxget",
	Short: "Show help for mxget commands",
	Long: `
 _____ ______      ___    ___ ________  _______  _________   
|\   _ \  _   \   |\  \  /  /|\   ____\|\  ___ \|\___   ___\ 
\ \  \\\__\ \  \  \ \  \/  / | \  \___|\ \   __/\|___ \  \_| 
 \ \  \\|__| \  \  \ \    / / \ \  \  __\ \  \_|/__  \ \  \  
  \ \  \    \ \  \  /     \/   \ \  \|\  \ \  \_|\ \  \ \  \ 
   \ \__\    \ \__\/  /\   \    \ \_______\ \_______\  \ \__\
    \|__|     \|__/__/ /\ __\    \|_______|\|_______|   \|__|
                  |__|/ \|__|                                

A simple tool that help you search and download your favorite music,
please visit https://github.com/winterssy/mxget for more detail.
`,
	Version: Version,
}

func Execute() error {
	return CmdRoot.Execute()
}

func init() {
	cobra.OnInitialize(settings.Init)

	CmdRoot.AddCommand(search.CmdSearch)
	CmdRoot.AddCommand(song.CmdSong)
	CmdRoot.AddCommand(artist.CmdArtist)
	CmdRoot.AddCommand(album.CmdAlbum)
	CmdRoot.AddCommand(playlist.CmdPlaylist)
	CmdRoot.AddCommand(serve.CmdServe)
	CmdRoot.AddCommand(config.CmdSet)
}
