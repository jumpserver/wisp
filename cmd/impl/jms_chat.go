package impl

import (
	"bytes"
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

func (j *JMServer) CallAPI(ctx context.Context, req *pb.HTTPRequest) (*pb.HTTPResponse, error) {
	method := req.GetMethod()
	reqUrl := req.GetPath()
	query := req.GetQuery()
	body := req.GetBody()
	header := req.GetHeader()
	apiClient := j.apiClient.Copy()
	var (
		status pb.Status
	)
	for k, v := range header {
		apiClient.SetHeader(k, v)
	}
	var bodyObj any
	if len(body) > 0 {
		if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&bodyObj); err != nil {
			logger.Errorf("Decode body failed: %v", err)
			status.Err = err.Error()
			return &pb.HTTPResponse{Status: &status}, err
		}
	}
	var resJson bytes.Buffer
	_, err := apiClient.Call(method, reqUrl, bodyObj, &resJson, query)
	if err != nil {
		logger.Errorf("Call: %v", err)
		status.Err = err.Error()
		return &pb.HTTPResponse{Status: &status}, err
	}
	status.Ok = true
	return &pb.HTTPResponse{Status: &status, Body: resJson.Bytes()}, nil
}
