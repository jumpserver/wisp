package cmd

import (
	"net"
	"path/filepath"

	"github.com/spf13/cobra"

	"github/jumpserver/wisp/cmd/common"
	"github/jumpserver/wisp/pkg/config"
	"github/jumpserver/wisp/pkg/logger"
)

var (
	cfgFile string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yml", "config file")
}

func initConfig() {
	config.Setup(cfgFile)
}

func initLogger() {
	conf := config.Get()
	logFilePath := filepath.Join(conf.LogDirPath, "wisp.log")
	logger.Setup(conf.LogLevel, logFilePath)
}

var rootCmd = &cobra.Command{
	Use:   "wisp",
	Short: "wisp is a grpc server to proxy JumpServer HTTP APIs",
	Long:  `A grpc server make easy to work with JumpServer HTTP APIs.`,
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		initLogger()
		conf := config.Get()
		ctx := common.GetSignalCtx()
		addr := net.JoinHostPort(conf.BindHost, conf.BindPort)
		srv := NewServer(addr, nil)
		Run(ctx, srv)
	}}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}
