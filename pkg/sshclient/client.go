package sshclient

import (
	"net"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
)

type Option func(conf *Options)

type Options struct {
	Host         string
	Port         string
	Username     string
	Password     string
	PrivateKey   string
	Passphrase   string
	keyboardAuth ssh.KeyboardInteractiveChallenge
	PrivateAuth  ssh.Signer
	timeout      int //  秒单位
}

func (cfg *Options) AuthMethods() []ssh.AuthMethod {
	authMethods := make([]ssh.AuthMethod, 0, 3)
	if cfg.PrivateAuth != nil {
		authMethods = append(authMethods, ssh.PublicKeys(cfg.PrivateAuth))
	}

	if cfg.PrivateKey != "" {
		var (
			signer ssh.Signer
			err    error
		)
		if cfg.Passphrase != "" {
			// 先使用 passphrase 解析 PrivateKey
			if signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(cfg.PrivateKey),
				[]byte(cfg.Passphrase)); err == nil {
				authMethods = append(authMethods, ssh.PublicKeys(signer))
			}
		}
		if err != nil || cfg.Passphrase == "" {
			// 1. 如果之前使用解析失败，则去掉 passphrase，则尝试直接解析 PrivateKey 防止错误的passphrase
			// 2. 如果没有 Passphrase 则直接解析 PrivateKey
			if signer, err = ssh.ParsePrivateKey([]byte(cfg.PrivateKey)); err == nil {
				authMethods = append(authMethods, ssh.PublicKeys(signer))
			}
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

func Username(username string) Option {
	return func(args *Options) {
		args.Username = username
	}
}

func Password(password string) Option {
	return func(args *Options) {
		args.Password = password
	}
}

func PrivateKey(privateKey string) Option {
	return func(args *Options) {
		args.PrivateKey = privateKey
	}
}

func Passphrase(passphrase string) Option {
	return func(args *Options) {
		args.Passphrase = passphrase
	}
}

func Host(host string) Option {
	return func(args *Options) {
		args.Host = host
	}
}

func Port(port int) Option {
	return func(args *Options) {
		args.Port = strconv.Itoa(port)
	}
}

func Timeout(timeout int) Option {
	return func(args *Options) {
		args.timeout = timeout
	}
}

func PrivateAuth(privateAuth ssh.Signer) Option {
	return func(args *Options) {
		args.PrivateAuth = privateAuth
	}
}

func KeyboardAuth(keyboardAuth ssh.KeyboardInteractiveChallenge) Option {
	return func(conf *Options) {
		conf.keyboardAuth = keyboardAuth
	}
}

const defaultTimeout = 15

func New(opts ...Option) (*ssh.Client, error) {
	cfg := &Options{
		Host:    "127.0.0.1",
		Port:    "22",
		timeout: defaultTimeout,
	}
	for _, setter := range opts {
		setter(cfg)
	}
	return NewWithCfg(cfg)
}

func NewWithCfg(cfg *Options) (*ssh.Client, error) {
	sshCfg := ssh.ClientConfig{
		User:            cfg.Username,
		Auth:            cfg.AuthMethods(),
		Timeout:         time.Duration(cfg.timeout) * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	destAddr := net.JoinHostPort(cfg.Host, cfg.Port)
	sshClient, err := ssh.Dial("tcp", destAddr, &sshCfg)
	if err != nil {
		return nil, err
	}
	return sshClient, nil
}
