package httpserver

import (
	"github.com/congcongke/httpfileserver/pkg/config"
	"github.com/congcongke/httpfileserver/pkg/server"

	"github.com/spf13/cobra"
)

func NewHttpFileServerCommand() *cobra.Command {
	conf := &config.Config{}
	cmd := &cobra.Command{
		Use:   "httpfileserver",
		Short: "httpfileserver is a simple file server through http",
		Long:  "it is expected to download file in compress mode via base auth",
		Run: func(cmd *cobra.Command, args []string) {
			server.Run(conf)
		},
	}

	cmd.PersistentFlags().StringVar(&conf.RootPath, "root", ".", "the root dir of file server")
	cmd.PersistentFlags().Uint16Var(&conf.Port, "port", 80, "the port exported outside")
	cmd.PersistentFlags().StringVar(&conf.Auth.Username, "user", "root", "username")
	cmd.PersistentFlags().StringVar(&conf.Auth.Password, "password", "", "user password")

	return cmd
}
