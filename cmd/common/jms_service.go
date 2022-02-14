package common

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github/jumpserver/wisp/pkg/config"
	"github/jumpserver/wisp/pkg/jms-sdk-go/model"
	"github/jumpserver/wisp/pkg/jms-sdk-go/service"
	"github/jumpserver/wisp/pkg/logger"
)

func MustJMService(conf config.Config) *service.JMService {
	key := MustLoadValidAccessKey(&conf)
	jmsService, err := service.NewAuthJMService(service.JMSCoreHost(
		conf.CoreHost), service.JMSTimeOut(30*time.Second),
		service.JMSAccessKey(key.ID, key.Secret),
	)
	if err != nil {
		logger.Fatal("创建 JMS Service 失败 " + err.Error())
		os.Exit(1)
	}
	return jmsService
}

func MustLoadValidAccessKey(conf *config.Config) model.AccessKey {
	var (
		coreHost       = conf.CoreHost
		name           = conf.Name
		bootstrapToken = conf.BootstrapToken
		dstFilePath    = conf.AccessKeyFilePath
	)
	if key, err := model.LoadAccessKeyFromFile(dstFilePath); err == nil {
		if err = MustValidKey(conf, key); err == nil {
			return key
		}
	}
	key := MustRegisterTerminalAccount(coreHost, name, bootstrapToken)
	if err := key.SaveToFile(dstFilePath); err != nil {
		logger.Error("保存key失败: " + err.Error())
	}
	return key
}

func MustRegisterTerminalAccount(coreHost, name, token string) (key model.AccessKey) {
	for i := 0; i < 10; i++ {
		terminal, err := service.RegisterTerminalAccount(coreHost, name, token)
		if err != nil {
			logger.Error(err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		key.ID = terminal.ServiceAccount.AccessKey.ID
		key.Secret = terminal.ServiceAccount.AccessKey.Secret
		return key
	}
	logger.Error("注册终端失败退出")
	os.Exit(1)
	return
}

func MustValidKey(conf *config.Config, key model.AccessKey) error {
	for i := 0; i < 10; i++ {
		if err := service.ValidAccessKey(conf.CoreHost, key); err != nil {
			switch {
			case errors.Is(err, service.ErrUnauthorized):
				logger.Error("Access key unauthorized, try to register new access key")
				return err
			default:
				logger.Error("校验 access key failed: " + err.Error())
			}
			time.Sleep(5 * time.Second)
			continue
		}
		return nil
	}
	logger.Error("校验 access key failed 退出")
	os.Exit(1)
	return fmt.Errorf("校验 access key %s 失败: ", conf.AccessKeyFilePath)
}
