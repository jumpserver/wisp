package cmd

import (
	"net"

	"github.com/spf13/cobra"

	"github.com/jumpserver/wisp/cmd/common"
	"github.com/jumpserver/wisp/cmd/impl"
	"github.com/jumpserver/wisp/pkg/config"
)

var fakeCmd = &cobra.Command{
	Use:   "fake",
	Short: "Run grpc server with test data",
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		conf := config.Get()
		ctx := common.GetSignalCtx()
		addr := net.JoinHostPort(conf.BindHost, conf.BindPort)
		implServer := impl.NewFakeServer(testFile)
		srv := NewServer(addr, implServer)
		RunServer(ctx, srv)
	},
}

var testFile string

func init() {
	rootCmd.AddCommand(fakeCmd)
	fakeCmd.Flags().StringVar(&testFile, "data", "test.json", "fake data")
}
