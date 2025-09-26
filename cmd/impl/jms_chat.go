package impl

import (
	"context"
	"encoding/json"

	"github.com/jumpserver/wisp/pkg/logger"
	pb "github.com/jumpserver/wisp/protobuf-go/protobuf"
	"google.golang.org/protobuf/types/known/structpb"
)

func (j *JMServer) GetAccountChat(ctx context.Context, req *pb.Empty) (*pb.AccountDetailResponse, error) {
	var (
		status pb.Status
	)
	accountChatDetail, err := j.apiClient.GetAccountsChat()
	if err != nil {
		logger.Errorf("GetAccountChat: %v", err)
		status.Err = err.Error()
		return &pb.AccountDetailResponse{
			Status: &status,
		}, nil
	}
	var m map[string]any
	jsonBytes, err := json.Marshal(accountChatDetail)
	if err != nil {
		logger.Errorf("Encode accountChatDetail failed: %v", err)
		status.Err = err.Error()
		return &pb.AccountDetailResponse{
			Status: &status,
		}, nil
	}
	if err = json.Unmarshal(jsonBytes, &m); err != nil {
		logger.Errorf("Unmarshal to map failed: %v", err)
		status.Err = err.Error()
		return &pb.AccountDetailResponse{
			Status: &status,
		}, nil
	}
	st, err1 := structpb.NewStruct(m)
	if err1 != nil {
		status.Err = err1.Error()
		return &pb.AccountDetailResponse{
			Status: &status,
		}, nil
	}
	status.Ok = true
	return &pb.AccountDetailResponse{Status: &status, Payload: st}, nil
}
