package sshclient

import (
	"golang.org/x/crypto/ssh"
	"strconv"
)

func WithUsername(username string) Option {
	return func(args *Config) {
		args.Username = username
	}
}

func WithPassword(password string) Option {
	return func(args *Config) {
		args.Password = password
	}
}

func WithPrivateKey(privateKey string) Option {
	return func(args *Config) {
		args.PrivateKey = privateKey
	}
}

func WithPassphrase(passphrase string) Option {
	return func(args *Config) {
		args.Passphrase = passphrase
	}
}

func WithHost(host string) Option {
	return func(args *Config) {
		args.Host = host
	}
}

func WithPort(port int) Option {
	return func(args *Config) {
		args.Port = strconv.Itoa(port)
	}
}

func WithTimeout(timeout int) Option {
	return func(args *Config) {
		args.Timeout = timeout
	}
}

func PrivateAuth(privateAuth ssh.Signer) Option {
	return func(args *Config) {
		args.PrivateAuth = privateAuth
	}
}

func KeyboardAuth(keyboardAuth ssh.KeyboardInteractiveChallenge) Option {
	return func(conf *Config) {
		conf.keyboardAuth = keyboardAuth
	}
}