package impl

import (
	"github/jumpserver/wisp/pkg/jms-sdk-go/service"
	pb "github/jumpserver/wisp/protobuf-go/protobuf"
)

type JmsService struct {
	pb.UnimplementedServiceServer
	apiClient *service.JMService
}

func (j *JmsService) GetAPIClient() *service.JMService {
	return j.apiClient
}
