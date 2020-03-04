package playlist

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/winterssy/glog"
	"github.com/winterssy/mxget/internal/cli"
	"github.com/winterssy/mxget/internal/settings"
	"github.com/winterssy/mxget/pkg/provider"
	"github.com/winterssy/mxget/pkg/utils"
)

var (
	playlistId string
	from       string
)

var CmdPlaylist = &cobra.Command{
	Use:   "playlist",
	Short: "Fetch and download playlist's songs via its id",
}

func Run(cmd *cobra.Command, args []string) {
	if playlistId == "" {
		playlistId = utils.Input("Playlist id")
		fmt.Println()
	}

	platform := settings.Cfg.Platform
	if from != "" {
		platform = from
	}

	client, err := provider.GetClient(platform)
	if err != nil {
		glog.Fatal(err)
	}

	glog.Infof("Fetch playlist [%s] from [%s]", playlistId, provider.GetDesc(platform))
	ctx := context.Background()
	playlist, err := client.GetPlaylist(ctx, playlistId)
	if err != nil {
		glog.Fatal(err)
	}

	cli.ConcurrentDownload(ctx, client, playlist.Name, playlist.Songs...)
}

func init() {
	CmdPlaylist.Flags().StringVar(&playlistId, "id", "", "playlist id")
	CmdPlaylist.Flags().StringVar(&from, "from", "", "music platform")
	CmdPlaylist.Flags().IntVar(&settings.Limit, "limit", 0, "concurrent download limit")
	CmdPlaylist.Flags().BoolVar(&settings.Tag, "tag", false, "update music metadata")
	CmdPlaylist.Flags().BoolVar(&settings.Lyric, "lyric", false, "download lyric")
	CmdPlaylist.Flags().BoolVar(&settings.Force, "force", false, "overwrite already downloaded music")
	CmdPlaylist.Run = Run
}
