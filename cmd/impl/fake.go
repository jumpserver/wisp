package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gofrs/uuid"

	"github.com/jumpserver/wisp/pkg/jms-sdk-go/model"
	"github.com/jumpserver/wisp/pkg/logger"
	pb "github.com/jumpserver/wisp/protobuf-go/protobuf"
)

func NewFakeServer(testFile string) *FakeServer {
	testData := NewTestData(testFile)
	return &FakeServer{
		testData: testData,
	}
}

type FakeServer struct {
	pb.UnimplementedServiceServer
	testData testData
}

func (f *FakeServer) GetDBTokenAuthInfo(ctx context.Context,
	req *pb.DBTokenRequest) (*pb.DBTokenResponse, error) {
	logger.Infof("Get DB Token AuthInfo req: %+v", req)
	var status pb.Status
	if f.testData.ID != req.Token {
		status.Err = fmt.Sprintf("not found %s", req.Token)
		return &pb.DBTokenResponse{Status: &status}, nil
	}
	status.Ok = true
	data := f.testData
	dbTokenInfo := pb.DBTokenAuthInfo{
		KeyId:       data.ID,
		SecreteId:   data.Secrete,
		Application: ConvertToProtobufApplication(data.Application),
		User:        ConvertToProtobufUser(data.User),
		SystemUser:  ConvertToProtobufSystemUser(data.SystemUserAuthInfo),
		Permission:  ConvertToProtobufPermission(data.Permission),
		ExpireInfo:  ConvertToProtobufExpireInfo(data.ExpireInfo),
		Gateways:    ConvertToProtobufGateWays(data.Gateways),
		FilterRules: ConvertToProtobufFilterRules(data.FilterRules),
	}

	return &pb.DBTokenResponse{Status: &status, Data: &dbTokenInfo}, nil
}

func (f *FakeServer) CreateSession(ctx context.Context,
	req *pb.SessionCreateRequest) (*pb.SessionCreateResponse, error) {
	logger.Infof("Create session req: %v", req)
	uid, _ := uuid.NewV4()
	status := pb.Status{Ok: true}
	return &pb.SessionCreateResponse{Status: &status,
		Data: &pb.Session{Id: uid.String()}}, nil
}

func (f *FakeServer) FinishSession(ctx context.Context,
	req *pb.SessionFinishRequest) (*pb.SessionFinishResp, error) {
	logger.Infof("Finish session req: %+v", req)
	status := pb.Status{Ok: true}
	return &pb.SessionFinishResp{Status: &status}, nil
}

func (f *FakeServer) UploadReplayFile(ctx context.Context,
	req *pb.ReplayRequest) (*pb.ReplayResponse, error) {
	logger.Infof("Upload replay file %+v", req)
	status := pb.Status{Ok: true}
	return &pb.ReplayResponse{
		Status: &status}, nil
}

func (f *FakeServer) UploadCommand(ctx context.Context,
	req *pb.CommandRequest) (*pb.CommandResponse, error) {
	logger.Infof("Upload Command %+v", req)
	status := pb.Status{Ok: true}
	return &pb.CommandResponse{
		Status: &status}, nil
}

func (f *FakeServer) DispatchTask(taskSrv pb.Service_DispatchTaskServer) error {
	for {
		req, err := taskSrv.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		logger.Infof("dispatch task %+v", req)
	}
}
func (f *FakeServer) ScanRemainReplays(ctx context.Context, req *pb.RemainReplayRequest) (*pb.RemainReplayResponse, error) {
	status := pb.Status{Ok: true}
	logger.Infof("Scan Remain Replays %+v", req)
	return &pb.RemainReplayResponse{Status: &status,
		SuccessFiles: []string{},
		FailureFiles: []string{},
		FailureErrs:  []string{},
	}, nil
}

type testData struct {
	ID                 string                   `json:"id"`
	Secrete            string                   `json:"secrete"`
	User               model.User               `json:"user"`
	Application        model.Application        `json:"application"`
	SystemUserAuthInfo model.SystemUserAuthInfo `json:"system_user_auth_info"`
	Gateways           []model.Gateway          `json:"gateways"`
	FilterRules        model.FilterRules        `json:"filter_rules"`
	ExpireInfo         model.ExpireInfo         `json:"expire_info"`
	Permission         model.Permission         `json:"permission"`
}

func NewTestData(filePath string) testData {
	var data testData
	fd, err := os.Open(filePath)
	if err != nil {
		logger.Fatal(err)
	}
	defer fd.Close()
	if err = json.NewDecoder(fd).Decode(&data); err != nil {
		logger.Fatal(err)
	}

	return data

}
