package common

import (
	"errors"

	"sync"

	"golang.org/x/crypto/ssh"

	"github.com/jumpserver/wisp/pkg/forward"
	"github.com/jumpserver/wisp/pkg/jms-sdk-go/model"
	"github.com/jumpserver/wisp/pkg/logger"
	"github.com/jumpserver/wisp/pkg/sshclient"

	pb "github.com/jumpserver/wisp/protobuf-go/protobuf"
)

func NewForwardCache() *ForwardCache {
	return &ForwardCache{
		lns: make(map[string]*forward.SSHForward),
	}
}

type ForwardCache struct {
	lns map[string]*forward.SSHForward

	sync.Mutex
}

func (s *ForwardCache) Add(key string, forward *forward.SSHForward) {
	s.Lock()
	defer s.Unlock()
	s.lns[key] = forward

}

func (s *ForwardCache) Remove(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.lns, key)

}

func (s *ForwardCache) Get(key string) *forward.SSHForward {
	s.Lock()
	defer s.Unlock()
	return s.lns[key]
}

var (
	ErrNoAvailable = errors.New("no available gateway")
)

func FindAvailableDomainGateway(domain *model.Domain) (*ssh.Client, error) {
	for i := range domain.Gateways {
		gateway := domain.Gateways[i]
		opts := make([]sshclient.Option, 0, 7)
		opts = append(opts, sshclient.WithHost(gateway.IP))
		opts = append(opts, sshclient.WithPort(gateway.Port))
		opts = append(opts, sshclient.WithUsername(gateway.Username))
		opts = append(opts, sshclient.WithPassword(gateway.Password))
		opts = append(opts, sshclient.WithPrivateKey(gateway.PrivateKey))
		opts = append(opts, sshclient.WithPassphrase(gateway.Password))
		opts = append(opts, sshclient.WithTimeout(15))
		proxyClient, err := sshclient.New(opts...)
		if err == nil {
			return proxyClient, nil
		}
		logger.Infof("Domain %s use gateway %s failed: %s",
			domain.Name, gateway.Name, err)
	}
	logger.Errorf("Domain %s find available gateway failed: %s", domain.Name)
	return nil, ErrNoAvailable
}

func FindAvailableGateway(gateways []*pb.Gateway) (*ssh.Client, error) {
	for i := range gateways {
		gateway := gateways[i]
		opts := make([]sshclient.Option, 0, 7)
		opts = append(opts, sshclient.WithHost(gateway.Ip))
		opts = append(opts, sshclient.WithPort(int(gateway.Port)))
		opts = append(opts, sshclient.WithUsername(gateway.Username))
		opts = append(opts, sshclient.WithPassword(gateway.Password))
		opts = append(opts, sshclient.WithPrivateKey(gateway.PrivateKey))
		opts = append(opts, sshclient.WithPassphrase(gateway.Password))
		opts = append(opts, sshclient.WithTimeout(15))
		proxyClient, err := sshclient.New(opts...)
		if err == nil {
			logger.Infof("Use gateway %s", gateway.Name)
			return proxyClient, nil
		}
	}
	return nil, ErrNoAvailable
}
