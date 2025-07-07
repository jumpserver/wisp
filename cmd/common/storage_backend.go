package common

import (
	"strings"

	"github.com/jumpserver-dev/sdk-go/model"
	"github.com/jumpserver-dev/sdk-go/service"
	"github.com/jumpserver/wisp/pkg/storage"
)

type StorageType interface {
	TypeName() string
}

type ReplayStorage interface {
	Upload(gZipFile, target string) error
	StorageType
}

type CommandStorage interface {
	BulkSave(commands []*model.Command) error
	StorageType
}

func NewCommandBackend(apiClient *service.JMService, cfg *model.CommandConfig) CommandStorage {
	switch cfg.TypeName {
	case "es", "elasticsearch":
		var hosts = cfg.Hosts
		var skipVerify bool
		index := cfg.Index
		docType := cfg.DocType
		if cfg.Other != nil {
			skipVerify = cfg.Other.IgnoreVerifyCerts
		}
		if index == "" {
			index = "jumpserver"
		}
		if docType == "" {
			docType = "_doc"
		}
		return storage.ESCommandStorage{
			Hosts:              hosts,
			Index:              index,
			DocType:            docType,
			InsecureSkipVerify: skipVerify,
		}
	case "null":
		return storage.NewNullStorage()
	default:
		return storage.ServerStorage{StorageType: "server", JmsService: apiClient}
	}
}

func NewReplayBackend(apiClient *service.JMService, cfg *model.ReplayConfig) ReplayStorage {

	switch cfg.TypeName {
	case "azure":
		var (
			accountName    string
			accountKey     string
			containerName  string
			endpointSuffix string
		)
		endpointSuffix = cfg.EndpointSuffix
		accountName = cfg.AccountName
		accountKey = cfg.AccountKey
		containerName = cfg.ContainerName

		if endpointSuffix == "" {
			endpointSuffix = "core.chinacloudapi.cn"
		}
		return storage.AzureReplayStorage{
			AccountName:    accountName,
			AccountKey:     accountKey,
			ContainerName:  containerName,
			EndpointSuffix: endpointSuffix,
		}
	case "oss":
		var (
			endpoint  string
			bucket    string
			accessKey string
			secretKey string
		)

		endpoint = cfg.Endpoint
		bucket = cfg.Bucket
		accessKey = cfg.AccessKey
		secretKey = cfg.SecretKey
		return storage.OSSReplayStorage{
			Endpoint:  endpoint,
			Bucket:    bucket,
			AccessKey: accessKey,
			SecretKey: secretKey,
		}
	case "s3", "swift", "cos":
		var (
			region    string
			endpoint  string
			bucket    string
			accessKey string
			secretKey string
		)

		bucket = cfg.Bucket
		endpoint = cfg.Endpoint
		region = cfg.Region
		accessKey = cfg.AccessKey
		secretKey = cfg.SecretKey

		if region == "" && endpoint != "" {
			endpointArray := strings.Split(endpoint, ".")
			if len(endpointArray) >= 2 {
				region = endpointArray[1]
			}
		}
		if bucket == "" {
			bucket = "jumpserver"
		}
		return storage.S3ReplayStorage{
			Bucket:    bucket,
			Region:    region,
			AccessKey: accessKey,
			SecretKey: secretKey,
			Endpoint:  endpoint,
		}
	case "obs":
		var (
			endpoint  string
			bucket    string
			accessKey string
			secretKey string
		)

		endpoint = cfg.Endpoint
		bucket = cfg.Bucket
		accessKey = cfg.AccessKey
		secretKey = cfg.SecretKey
		return storage.OBSReplayStorage{
			Endpoint:  endpoint,
			Bucket:    bucket,
			AccessKey: accessKey,
			SecretKey: secretKey,
		}
	case "null":
		return storage.NewNullStorage()
	default:
		return storage.ServerStorage{StorageType: "server", JmsService: apiClient}
	}
}
