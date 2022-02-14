package impl

import (
	"github.com/jumpserver/wisp/pkg/jms-sdk-go/service"
	pb "github.com/jumpserver/wisp/protobuf-go/protobuf"
)

type JmsService struct {
	pb.UnimplementedServiceServer
	apiClient *service.JMService
}

func (j *JmsService) GetAPIClient() *service.JMService {
	return j.apiClient
}
