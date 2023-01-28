package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/EtoNeJa00/GRPCService/internal/app/mock"
	"github.com/EtoNeJa00/GRPCService/internal/transport/grpc/generated/record"
)

type grpcServerTestSuite struct {
	suite.Suite
	port    string
	mockEnp *mock.MockGrpcEnp
	grpcS   GRPCServer
	ctx     context.Context
}

func (g *grpcServerTestSuite) SetupSuite() {
	g.ctx = context.Background()
	g.port = ":8080"

	ctrl := gomock.NewController(g.T())
	g.mockEnp = mock.NewMockGrpcEnp(ctrl)

	var err error

	g.grpcS, err = NewGRPCServer(g.port, g.mockEnp)
	g.Require().NoError(err)
	go func() {
		err := g.grpcS.StartServer()
		g.Require().NoError(err)
	}()
}

func (g *grpcServerTestSuite) TearDownSuite() {
	if g.grpcS != nil {
		g.grpcS.StopServer()
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(grpcServerTestSuite))
}

func (g *grpcServerTestSuite) TestGrpc() {
	conn, err := grpc.Dial("127.0.0.1"+g.port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	g.Require().NoError(err)

	defer func() {
		err := conn.Close()
		g.Require().NoError(err)
	}()

	id := uuid.New()
	data := "data"
	client := record.NewRecordsClient(conn)
	request := record.Record{
		Id:   uuid.Nil.String(),
		Data: data,
	}

	g.testSet(data, id, client, &request)

	g.testGet(id, data, client)

	g.testDelete(id, data, err, client)
}

func (g *grpcServerTestSuite) testSet(data string, id uuid.UUID, client record.RecordsClient, request *record.Record) {
	g.mockEnp.EXPECT().Set(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, r *record.Record) (*record.Record, error) {
		if r.GetId() != uuid.Nil.String() || r.GetData() != data {
			return nil, errors.New("")
		}

		return &record.Record{
			Id:   id.String(),
			Data: data,
		}, nil
	})

	response, err := client.Set(g.ctx, request)
	g.Require().NoError(err)
	g.Require().EqualValues(id.String(), response.Id)
}

func (g *grpcServerTestSuite) testGet(id uuid.UUID, data string, client record.RecordsClient) {
	g.mockEnp.EXPECT().Get(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, r *record.Id) (*record.Record, error) {
		if r.GetId() != id.String() {
			return nil, errors.New("")
		}

		return &record.Record{
			Id:   id.String(),
			Data: data,
		}, nil
	})

	reqId := record.Id{Id: id.String()}

	response, err := client.Get(g.ctx, &reqId)
	g.Require().NoError(err)
	g.Require().EqualValues(id.String(), response.Id)
}

func (g *grpcServerTestSuite) testDelete(id uuid.UUID, data string, err error, client record.RecordsClient) {
	g.mockEnp.EXPECT().Delete(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, r *record.Id) (*record.Record, error) {
		if r.GetId() != id.String() {
			return nil, errors.New("")
		}

		return &record.Record{
			Id:   id.String(),
			Data: data,
		}, nil
	})

	reqId := record.Id{Id: id.String()}

	response, err := client.Delete(g.ctx, &reqId)
	g.Require().NoError(err)
	g.Require().EqualValues(id.String(), response.Id)
}
