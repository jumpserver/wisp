package cmd

import (
	"context"
	"net"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/jumpserver/wisp/cmd/common"
	"github.com/jumpserver/wisp/cmd/impl"
	"github.com/jumpserver/wisp/pkg/config"
	"github.com/jumpserver/wisp/pkg/logger"
	"github.com/jumpserver/wisp/pkg/process"
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
	logFilePath := filepath.Join(conf.LogFolderPath, "wisp.log")
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
		addr := net.JoinHostPort(conf.BindHost, conf.BindPort)
		apiClient := common.MustJMService(&conf)
		cfg, err := apiClient.GetTerminalConfig()
		if err != nil {
			logger.Fatalf("Get Terminal Cfg failed: %s", err)
		}
		uploader := common.NewUploader(apiClient, &cfg)
		{
			go uploader.Start()
			defer uploader.Stop()
		}
		beat := common.NewBeatService(apiClient)
		{
			go beat.KeepHeartBeat()
		}
		ctx := common.GetSignalCtx()
		grpcImplSrv := impl.NewJMServer(apiClient, uploader, beat)
		srv := NewServer(addr, grpcImplSrv)
		{
			go RunServer(ctx, srv)
			time.Sleep(time.Second)
		}
		{
			//  子进程启动
			if conf.ExecuteProgram != "" {
				subProcess := process.New(conf.RootPath, conf.ExecuteProgram)
				logger.Infof("start subprocess: %s", conf.ExecuteProgram)
				startProcess(ctx, subProcess)
				// 子进程突出，srv 结束
				srv.Stop()
			}
		}
		srv.Wait()
	}}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}

func startProcess(ctx context.Context, subProcess *process.Process) {
	go func() {
		<-ctx.Done()
		subProcess.Stop()
	}()
	if err := subProcess.Start(); err != nil {
		logger.Fatal(err)
	}
}
