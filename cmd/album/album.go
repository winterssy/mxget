package album

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
	albumId string
	from    string
)

var CmdAlbum = &cobra.Command{
	Use:   "album",
	Short: "Fetch and download album's songs via its id",
}

func Run(cmd *cobra.Command, args []string) {
	if albumId == "" {
		albumId = utils.Input("Album id")
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

	glog.Infof("Fetch album [%s] from [%s]", albumId, provider.GetDesc(platform))
	ctx := context.Background()
	album, err := client.GetAlbum(ctx, albumId)
	if err != nil {
		glog.Fatal(err)
	}

	cli.ConcurrentDownload(ctx, client, album.Name, album.Songs...)
}

func init() {
	CmdAlbum.Flags().StringVar(&albumId, "id", "", "album id")
	CmdAlbum.Flags().StringVar(&from, "from", "", "music platform")
	CmdAlbum.Flags().IntVar(&settings.Limit, "limit", 0, "concurrent download limit")
	CmdAlbum.Flags().BoolVar(&settings.Tag, "tag", false, "update music metadata")
	CmdAlbum.Flags().BoolVar(&settings.Lyric, "lyric", false, "download lyric")
	CmdAlbum.Flags().BoolVar(&settings.Force, "force", false, "overwrite already downloaded music")
	CmdAlbum.Run = Run
}
