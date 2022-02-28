package common

import (
	"errors"
	"os"
	"time"

	"github.com/jumpserver/wisp/pkg/config"
	"github.com/jumpserver/wisp/pkg/jms-sdk-go/model"
	"github.com/jumpserver/wisp/pkg/jms-sdk-go/service"
	"github.com/jumpserver/wisp/pkg/logger"
)

func MustJMService(conf *config.Config) *service.JMService {
	key := MustLoadValidAccessKey(conf)
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
		dstFilePath = conf.AccessKeyFilePath
	)
	if key, err := model.LoadAccessKeyFromFile(dstFilePath); err == nil {
		MustValidKey(conf, &key)
		return key
	}
	return MustRegisterTerminal(conf)
}

func MustRegisterTerminal(conf *config.Config) (key model.AccessKey) {
	var (
		componentName  = conf.ComponentName
		coreHost       = conf.CoreHost
		name           = conf.Name
		bootstrapToken = conf.BootstrapToken
		dstFilePath    = conf.AccessKeyFilePath
	)
	if _, ok := model.SupportedComponent(componentName); !ok {
		logger.Fatalf("组件名称错误: %s", componentName)
	}

	for i := 0; i < 10; i++ {
		terminal, err := service.RegisterTerminalAccount(coreHost, componentName, name, bootstrapToken)
		if err != nil {
			logger.Error(err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		key.ID = terminal.ServiceAccount.AccessKey.ID
		key.Secret = terminal.ServiceAccount.AccessKey.Secret
		if err2 := key.SaveToFile(dstFilePath); err2 != nil {
			logger.Errorf("保存key失败: %s", err)
		}
		return key
	}
	logger.Fatal("注册终端失败退出")
	return
}

func MustValidKey(conf *config.Config, key *model.AccessKey) {
	for i := 0; i < 10; i++ {
		if err := service.ValidAccessKey(conf.CoreHost, *key); err != nil {
			switch {
			case errors.Is(err, service.ErrUnauthorized):
				logger.Error("Access key 已失效, 重新注册")
				newKey := MustRegisterTerminal(conf)
				key.ID = newKey.ID
				key.Secret = newKey.Secret
			default:
				logger.Error("校验 access key failed: " + err.Error())
			}
			time.Sleep(5 * time.Second)
			continue
		}
		return
	}
	logger.Fatal("校验 access key failed 退出")
}
