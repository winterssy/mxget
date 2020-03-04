package serve

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/winterssy/mxget/internal/server"
)

var (
	port int
)

var CmdServe = &cobra.Command{
	Use:   "serve",
	Short: "Run mxget as an API server",
}

func Run(cmd *cobra.Command, args []string) {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	server.Init(app)
	app.Run(fmt.Sprintf(":%d", port))
}

func init() {
	CmdServe.Flags().IntVar(&port, "port", 8080, "server listening port")
	CmdServe.Run = Run
}
