package song

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
	songId string
	from   string
)

var CmdSong = &cobra.Command{
	Use:   "song",
	Short: "Fetch and download single song via its id",
}

func Run(cmd *cobra.Command, args []string) {
	if songId == "" {
		songId = utils.Input("Song id")
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

	glog.Infof("Fetch song [%s] from [%s]", songId, provider.GetDesc(platform))
	ctx := context.Background()
	song, err := client.GetSong(ctx, songId)
	if err != nil {
		glog.Fatal(err)
	}

	cli.ConcurrentDownload(ctx, client, ".", song)
}

func init() {
	CmdSong.Flags().StringVar(&songId, "id", "", "song id")
	CmdSong.Flags().StringVar(&from, "from", "", "music platform")
	CmdSong.Flags().BoolVar(&settings.Tag, "tag", false, "update music metadata")
	CmdSong.Flags().BoolVar(&settings.Lyric, "lyric", false, "download lyric")
	CmdSong.Flags().BoolVar(&settings.Force, "force", false, "overwrite already downloaded music")
	CmdSong.Run = Run
}
