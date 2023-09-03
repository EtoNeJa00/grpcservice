package endpoints

//go:generate go run -mod=mod github.com/golang/mock/mockgen -package=mock -destination=../app/mock/grpc_endpoint_generated.go -build_flags=-mod=mod . GrpcEnp

import (
	"context"

	"GRPCService/internal/app/usecase"
	"GRPCService/internal/models"
	"GRPCService/internal/transport/grpc/generated/record"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcEnp interface {
	Get(context.Context, *record.Id) (*record.Record, error)
	Set(context.Context, *record.Record) (*record.Record, error)
	Delete(context.Context, *record.Id) (*record.Record, error)
}

type grpcEnp struct {
	ctx context.Context
	uc  usecase.Usecase
}

func NewEndpoint(ctx context.Context, uc usecase.Usecase) GrpcEnp {
	return grpcEnp{ctx: ctx, uc: uc}
}

func (e grpcEnp) Get(ctx context.Context, id *record.Id) (*record.Record, error) {
	idUUID, err := uuid.Parse(id.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	rec, err := e.uc.GetRecord(ctx, idUUID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return e.ToPBRecord(rec), nil
}

func (e grpcEnp) Set(ctx context.Context, r *record.Record) (*record.Record, error) {
	rec, err := e.FromPBRecord(r)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	newRec, err := e.uc.SetRecord(ctx, rec)
	if err != nil {
		return nil, err
	}

	return e.ToPBRecord(newRec), nil
}

func (e grpcEnp) Delete(ctx context.Context, id *record.Id) (*record.Record, error) {
	idUUID, err := uuid.Parse(id.GetId())
	if err != nil {
		return nil, err
	}

	rec, err := e.uc.DeleteRecord(ctx, idUUID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return e.ToPBRecord(rec), nil
}

func (e grpcEnp) ToPBRecord(r models.Record) *record.Record {
	return &record.Record{
		Id:   r.ID.String(),
		Data: r.Data,
	}
}

func (e grpcEnp) FromPBRecord(r *record.Record) (models.Record, error) {
	id, err := uuid.Parse(r.GetId())
	if err != nil {
		return models.Record{}, status.Error(codes.InvalidArgument, err.Error())
	}

	return models.Record{
		ID:   id,
		Data: r.GetData(),
	}, err
}
