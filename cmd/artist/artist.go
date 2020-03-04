package artist

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
	artistId string
	from     string
)

var CmdArtist = &cobra.Command{
	Use:   "artist",
	Short: "Fetch and download artist's hot songs via its id",
}

func Run(cmd *cobra.Command, args []string) {
	if artistId == "" {
		artistId = utils.Input("Artist id")
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

	glog.Infof("Fetch artist [%s] from [%s]", artistId, provider.GetDesc(platform))
	ctx := context.Background()
	artist, err := client.GetArtist(ctx, artistId)
	if err != nil {
		glog.Fatal(err)
	}

	cli.ConcurrentDownload(ctx, client, artist.Name, artist.Songs...)
}

func init() {
	CmdArtist.Flags().StringVar(&artistId, "id", "", "artist id")
	CmdArtist.Flags().StringVar(&from, "from", "", "music platform")
	CmdArtist.Flags().IntVar(&settings.Limit, "limit", 0, "concurrent download limit")
	CmdArtist.Flags().BoolVar(&settings.Tag, "tag", false, "update music metadata")
	CmdArtist.Flags().BoolVar(&settings.Lyric, "lyric", false, "download lyric")
	CmdArtist.Flags().BoolVar(&settings.Force, "force", false, "overwrite already downloaded music")
	CmdArtist.Run = Run
}
