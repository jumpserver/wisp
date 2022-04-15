package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	ComponentName  string `mapstructure:"COMPONENT_NAME"`
	Name           string `mapstructure:"NAME"`
	CoreHost       string `mapstructure:"CORE_HOST"`
	BootstrapToken string `mapstructure:"BOOTSTRAP_TOKEN"`
	BindHost       string `mapstructure:"BIND_HOST"`
	BindPort       string `mapstructure:"BIND_PORT"`
	LogLevel       string `mapstructure:"LOG_LEVEL"`
	RootPath       string `mapstructure:"WORK_DIR"`
	ExecuteProgram string `mapstructure:"EXECUTE_PROGRAM"`

	DataFolderPath    string
	LogFolderPath     string
	KeyFolderPath     string
	AccessKeyFilePath string
}

func (c *Config) ensureFolders() {
	if c.Name == "" {
		c.Name = getDefaultNameByComponent(c.ComponentName)
	}
	if c.RootPath == "" {
		c.RootPath = getPwdDirPath()
	}
	dataFolderPath := filepath.Join(c.RootPath, "data")
	LogFolderPath := filepath.Join(dataFolderPath, "logs")
	keyFolderPath := filepath.Join(dataFolderPath, "keys")
	accessKeyFilePath := filepath.Join(keyFolderPath, ".access_key")

	folders := []string{dataFolderPath, LogFolderPath, keyFolderPath}
	for i := range folders {
		if err := EnsureDirExist(folders[i]); err != nil {
			log.Fatalf("Create folder failed: %s\n", err)
		}
	}
	c.DataFolderPath = dataFolderPath
	c.LogFolderPath = LogFolderPath
	c.KeyFolderPath = keyFolderPath
	c.AccessKeyFilePath = accessKeyFilePath
}

var (
	GlobalConfig *Config
	once         sync.Once
)

func Get() Config {
	once.Do(func() {
		if GlobalConfig == nil {
			conf := getDefaultConfig()
			conf.ensureFolders()
			GlobalConfig = &conf
		}
	})
	return *GlobalConfig
}

func Setup(configPath string) {
	var conf = getDefaultConfig()
	loadConfigFromEnv(&conf)
	loadConfigFromFile(configPath, &conf)
	conf.ensureFolders()
	GlobalConfig = &conf
	fmt.Printf("%+v\n", conf)
}

func getDefaultConfig() Config {
	return Config{
		CoreHost:       "http://localhost:8080",
		BootstrapToken: "",
		BindHost:       "0.0.0.0",
		BindPort:       "9090",
		LogLevel:       "INFO",
	}

}

const (
	hostEnvKey = "SERVER_HOSTNAME"

	defaultNameMaxLen = 128
)

func getDefaultNameByComponent(component string) string {
	hostname, _ := os.Hostname()
	if serverHostname, ok := os.LookupEnv(hostEnvKey); ok {
		hostname = fmt.Sprintf("%s-%s", serverHostname, hostname)
	}
	prefixName := fmt.Sprintf("[%s]-", strings.Title(component))
	hostRune := []rune(prefixName + hostname)
	if len(hostRune) <= defaultNameMaxLen {
		return string(hostRune)
	}
	name := make([]rune, defaultNameMaxLen)
	index := defaultNameMaxLen / 2
	copy(name[:index], hostRune[:index])
	start := len(hostRune) - index
	copy(name[index:], hostRune[start:])
	return string(name)
}

func getPwdDirPath() string {
	if rootPath, err := os.Getwd(); err == nil {
		return rootPath
	}
	return ""
}

func loadConfigFromEnv(conf *Config) {
	viper.AutomaticEnv() // 全局配置，用于其他 pkg 包可以用 viper 获取环境变量的值
	envViper := viper.New()
	for _, item := range os.Environ() {
		envItem := strings.SplitN(item, "=", 2)
		if len(envItem) == 2 {
			envViper.Set(envItem[0], viper.Get(envItem[0]))
		}
	}
	if err := envViper.Unmarshal(conf); err == nil {
		fmt.Println("Load config from env")
	}
}

func loadConfigFromFile(path string, conf *Config) {
	var err error
	if have(path) {
		fileViper := viper.New()
		fileViper.SetConfigFile(path)
		if err = fileViper.ReadInConfig(); err == nil {
			if err = fileViper.Unmarshal(conf); err == nil {
				fmt.Printf("Load config from %s success\n", path)
				return
			}
		}
	}
	if err != nil {
		fmt.Printf("Load config from %s failed: %s\n", path, err)
		os.Exit(1)
	}
}

func have(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
func haveDir(file string) bool {
	fi, err := os.Stat(file)
	return err == nil && fi.IsDir()
}

func EnsureDirExist(path string) error {
	if !haveDir(path) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
