package sshclient

import (
	"errors"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

type Option func(conf *Config)

type Config struct {
	Host         string
	Port         string
	Username     string
	Password     string
	PrivateKey   string
	Passphrase   string
	keyboardAuth ssh.KeyboardInteractiveChallenge
	PrivateAuth  ssh.Signer
	Timeout      int //  秒单位
}

func (cfg *Config) AuthMethods() []ssh.AuthMethod {
	authMethods := make([]ssh.AuthMethod, 0, 3)
	if cfg.PrivateAuth != nil {
		authMethods = append(authMethods, ssh.PublicKeys(cfg.PrivateAuth))
	}
	if cfg.PrivateKey != "" {
		if signer, err := ParsePrivateKey(cfg.PrivateKey, cfg.Passphrase); err == nil {
			authMethods = append(authMethods, ssh.PublicKeys(signer))
		}
	}
	if cfg.Password != "" {
		authMethods = append(authMethods, ssh.Password(cfg.Password))
	}
	if cfg.keyboardAuth == nil {
		cfg.keyboardAuth = func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
			if len(questions) == 0 {
				return []string{}, nil
			}
			return []string{cfg.Password}, nil
		}
	}
	authMethods = append(authMethods, ssh.KeyboardInteractive(cfg.keyboardAuth))

	return authMethods
}

const defaultTimeout = 15

func New(opts ...Option) (*ssh.Client, error) {
	cfg := &Config{
		Host:    "127.0.0.1",
		Port:    "22",
		Timeout: defaultTimeout,
	}
	for _, setter := range opts {
		setter(cfg)
	}
	return NewWithCfg(cfg)
}

func NewWithCfg(cfg *Config) (*ssh.Client, error) {
	sshCfg := ssh.ClientConfig{
		User:            cfg.Username,
		Auth:            cfg.AuthMethods(),
		Timeout:         time.Duration(cfg.Timeout) * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	destAddr := net.JoinHostPort(cfg.Host, cfg.Port)
	sshClient, err := ssh.Dial("tcp", destAddr, &sshCfg)
	if err != nil {
		return nil, err
	}
	return sshClient, nil
}

var ErrNoAvailable = errors.New("no available config ")

func GetFirstAvailableClient(cfgs ...Config) (*ssh.Client, error) {
	for i := range cfgs {
		if client, err := NewWithCfg(&cfgs[i]); err == nil {
			return client, nil
		}
	}
	return nil, ErrNoAvailable
}

func ParsePrivateKey(privateKey, passphrase string) (signer ssh.Signer, err error) {
	if passphrase != "" {
		if signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(privateKey),
			[]byte(passphrase)); err == nil {
			return signer, nil
		}
	}
	// 1. 如果之前使用解析失败，则去掉 passphrase，则尝试直接解析 PrivateKey 防止错误的passphrase
	// 2. 如果没有 Passphrase 则直接解析 PrivateKey
	if signer, err = ssh.ParsePrivateKey([]byte(privateKey)); err == nil {
		return signer, nil
	}
	return nil, err
}
