package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/winterssy/glog"
	"github.com/winterssy/mxget/internal/settings"
	"github.com/winterssy/mxget/pkg/provider"
)

var (
	dir   string
	from  string
	show  bool
	reset bool
)

var CmdSet = &cobra.Command{
	Use:   "config",
	Short: "Specify the default behavior of mxget",
}

func Run(cmd *cobra.Command, args []string) {
	if cmd.Flags().NFlag() == 0 {
		_ = cmd.Help()
		return
	}

	if show {
		fmt.Print(fmt.Sprintf(`
      download dir -> %s
    music platform -> %s [%s]
`, settings.Cfg.Dir, settings.Cfg.Platform, provider.GetDesc(settings.Cfg.Platform)))
		return
	}

	if reset {
		settings.Cfg.Reset()
		return
	}

	if dir != "" {
		dir = filepath.Clean(dir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			glog.Fatalf("Can't make download dir: %v", err)
		}
		settings.Cfg.Dir = dir
	}
	if from != "" {
		if provider.GetDesc(from) == "unknown" {
			glog.Fatalf("Unexpected music platform: %q", from)
		}
		settings.Cfg.Platform = from
	}

	if dir != "" || from != "" {
		_ = settings.Cfg.Save()
	}
}

func init() {
	CmdSet.Flags().StringVar(&dir, "dir", "", "specify the default download directory")
	CmdSet.Flags().StringVar(&from, "from", "", "specify the default music platform")
	CmdSet.Flags().BoolVar(&show, "show", false, "show current settings")
	CmdSet.Flags().BoolVar(&reset, "reset", false, "reset default settings")
	CmdSet.Run = Run
}
